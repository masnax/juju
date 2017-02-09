#!/usr/bin/env python
"""Assess network health for a given deployment or bundle"""
from __future__ import print_function

import argparse
import logging
import sys
import json

from deploy_stack import (
    BootstrapManager,
    )
from utility import (
    add_basic_testing_arguments,
    configure_logging,
    )


__metaclass__ = type

log = logging.getLogger("assess_network_health")


def assess_network_health(client, bundle=None):
    # If a bundle is supplied, deploy it
    if bundle:
        client.deploy_bundle(bundle)
    # Else deploy two dummy charms to test on
    else:
        client.deploy('ubuntu', num=2)
    client.deploy('/trusty/network-health')
    # Wait for the deployment to finish.
    client.wait_for_started()
    # Grab apps from status
    apps = client.get_status()['applications']
    for service in apps.keys():
        self.client.juju('add-relation', service, 'network-health')
    log.info("Starting network tests")
    targets = parse_targets(apps)
    log.info(neighbor_visibility(client, targets))
    # Expose dummy charm if no bundle is specified
    if bundle is None:
        client.juju('expose', ('ubuntu',))
    # Grab exposed charms
    exposed = [app for app, info in apps.items() if info['exposed'] is True]
    # If we have exposed charms, test their exposure
    if len(exposed) > 0:
        log.info(ensure_exposed(client, targets, exposed))


def neighbor_visibility(client, targets):
    """Check if each application's units are visible, including our own.
    :param targets: Dict of units & public-addresses by application
    """
    results = {}
    # For each unit of network health, try to ping all known nodes
    for nh_unit in apps['network-health']['units']:
        service_resuts = {}
        for service, units in targets.items():
            # Change our dictionary into a json string
            units = json.dumps(units, separators=(',', '='))
            # Replace curly brackets so juju doesn't puke
            units.replace('{', '(')
            units.replace('}', ')')
            retval = new_client.action_do_fetch(unit, 'ping',
                                                'targets={}'.format(units))
            service_results[service] = retval
        results[nh_unit] = service_results
    return results


def ensure_exposed(client, targets, exposed):
    """Ensure exposed services are visible from the outside
    :param targets: Dict of units & public-addresses by application
    :param exposed: List of exposed services
    """
    # Spin up new client and deploy under it
    new_client = client.add_model(client.env)
    new_client.deploy('local:trusty/ubuntu')
    new_client.deploy('local:trusty/network-health')
    new_client.juju('add-relation' ('ubuntu', 'network-health'))
    # For each service, try to ping it from the outside model.
    service_results = {}
    for service, units in targets.items():
        units = json.dumps(units, separators=(',', '='))
        units.replace('{', '(')
        units.replace('}', ')')
        retval = new_client.action_do_fetch('network-health/0', 'ping',
                                            'targets={}'.format(units))
        service_results[service] = retval
    # Check revtal against exposed, return passes & failures
    fails = []
    passes = []
    for service, returns in service_results:
        if True in returns and service not in exposed:
            fails.append(service)
        elif True in returns and service in exposed:
            passes.append(service)
    return passes, failures


def parse_targets(apps):
    """Returns targets based on supplied juju status information.
    :param apps: Dict of applications via 'juju status --format yaml'
    """
    targets = {}
    for app, units in dic.items():
        target_units = {}
        for unit_id, info in units.get('units').items():
            target_units[unit_id] = info.get('public-address')
        targets[app] = target_units
    return targets


def parse_args(argv):
    """Parse all arguments."""
    parser = argparse.ArgumentParser(description="Test Network Health")
    add_basic_testing_arguments(parser)
    parser.add_argument('--bundle', help='Bundle to test network against')
    parser.set_defaults(series='trusty')
    return parser.parse_args(argv)


def main(argv=None):
    args = parse_args(argv)
    configure_logging(args.verbose)
    bs_manager = BootstrapManager.from_args(args)
    with bs_manager.booted_context(args.upload_tools):
        bundle = args.bundle
        assess_network_health(bs_manager.client, bundle)
    return 0


if __name__ == '__main__':
    sys.exit(main())
