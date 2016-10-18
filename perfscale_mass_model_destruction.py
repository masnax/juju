#!/usr/bin/env python
"""Perfscale test measuring adding and destroying a large number of models.

Steps taken in this test:
  - Bootstraps a provider
  - Creates x amount of models and waits for them to be ready
  - Delete all the models at once.
"""

import argparse
from datetime import datetime
import logging
import sys

from deploy_stack import (
    BootstrapManager,
)
from generate_perfscale_results import (
    DeployDetails,
    TimingData,
    run_perfscale_test,
)
from utility import (
    add_basic_testing_arguments,
    configure_logging,
)


log = logging.getLogger("perfscale_mass_model_destruction")

__metaclass__ = type


def perfscale_assess_model_destruction(client, args):
    """Create a bunch of models and then destroy them all."""
    model_count = args.model_count

    all_models = []
    for item in xrange(0, model_count):
        model_name = 'model{}'.format(item)
        log.info('Creating model: {}'.format(model_name))
        new_model = client.add_model(client.env.clone(model_name))
        new_model.wait_for_started()
        all_models.append(new_model)

    destruction_start = datetime.utcnow()
    for doomed in all_models:
        doomed.destroy_model()
    destruction_end = datetime.utcnow()

    destruction_timing = TimingData(destruction_start, destruction_end)
    return DeployDetails(
        'Destroy {} models'.format(model_count),
        {'Model Count': model_count},
        destruction_timing)


def parse_args(argv):
    """Parse all arguments."""
    parser = argparse.ArgumentParser(
        description="Perfscale bundle deployment test.")
    add_basic_testing_arguments(parser)
    parser.add_argument(
        '--model-count',
        type=int,
        help='Number of models to create.',
        default=100)
    return parser.parse_args(argv)


def main(argv=None):
    args = parse_args(argv)
    configure_logging(args.verbose)
    bs_manager = BootstrapManager.from_args(args)
    run_perfscale_test(perfscale_assess_model_destruction, bs_manager, args)

    return 0


if __name__ == '__main__':
    sys.exit(main())
