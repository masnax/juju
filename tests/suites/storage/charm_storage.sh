assess_rootfs() {
	echo "Assessing filesystem rootfs"
	# Assess charm storage with the filesystem storage provider
	juju deploy -m "${model_name}" ./testcharms/charms/dummy-storage-fs --series jammy --storage data=rootfs,1G
	if [ $(unit_exist "data/0") ]; then
		return 0
	fi
	wait_for_storage "attached" '.storage["data/0"]["status"].current'
	# assert the storage kind name
	assert_storage "filesystem" "$(kind_name "data" 0)"
	# assert the storage status
	assert_storage "alive" "$(life_status "data" 0)"
	# assert the storage label name
	assert_storage "data/0" "$(label 0)"
	# assert the unit attachment name
	assert_storage "dummy-storage-fs/0" "$(unit_attachment "data" 0 0)"
	# assert the attached unit state
	assert_storage "alive" "$(unit_state "data" 0 "dummy-storage-fs" 0)"
	wait_for_storage "attached" "$(filesystem_status 0 0).current"
	# assert the filesystem size
	requested_storage=1024
	acquired_storage=$(juju storage --format json | jq '.filesystems | .["0/0"] | select(.pool=="rootfs") | .size ')
	if [ "$requested_storage" -gt "$acquired_storage" ]; then
		echo "acquired storage size $acquired_storage should be greater than the requested storage $requested_storage"
		exit 1
	fi
	echo "Filesystem rootfs PASSED"
}

assess_loop_disk() {
	# Assess charm storage with the filesystem storage provider
	echo "Assessing block loop disk 1"
	juju deploy -m "${model_name}" ./testcharms/charms/dummy-storage-lp --series jammy --storage disks=loop,1G
	if [ $(unit_exist "disks/1") ]; then
		return 0
	fi
	wait_for_storage "attached" '.storage["disks/1"]["status"].current'
	# assert the storage kind name
	assert_storage "block" "$(kind_name "disks" 1)"
	# assert the storage status
	assert_storage "alive" "$(life_status "disks" 1)"
	# assert the storage label name
	assert_storage "disks/1" "$(label 1)"
	# assert the unit attachment name
	assert_storage "dummy-storage-lp/0" "$(unit_attachment "disks" 1 0)"
	# assert the attached unit state
	assert_storage "alive" "$(unit_state "disks" 1 "dummy-storage-lp" 0)"
	echo "Block loop disk 1 PASSED"

	echo "Assessing add storage block loop disk 2"
	juju add-storage -m "${model_name}" dummy-storage-lp/0 disks=1
	if [ $(unit_exist "disks/2") ]; then
		return 0
	fi
	wait_for_storage "attached" '.storage["disks/2"]["status"].current'
	# assert the storage kind name
	assert_storage "block" "$(kind_name "disks" 2)"
	# assert the storage status
	assert_storage "alive" "$(life_status "disks" 2)"
	# assert the storage label name
	assert_storage "disks/2" "$(label 2)"
	# assert the unit attachment name
	assert_storage "dummy-storage-lp/0" "$(unit_attachment "disks" 2 0)"
	# assert the attached unit state
	assert_storage "alive" "$(unit_state "disks" 2 "dummy-storage-lp" 0)"
	echo "Block loop disk 2 PASSED"
}

assess_tmpfs() {
	echo "Assessing filesystem tmpfs"
	juju deploy -m "${model_name}" ./testcharms/charms/dummy-storage-fs dummy-storage-tp --series jammy --storage data=tmpfs,1G
	if [ $(unit_exist "data/3") ]; then
		return 0
	fi
	wait_for_storage "attached" '.storage["data/3"]["status"].current'
	# assert the storage kind name
	assert_storage "filesystem" "$(kind_name "data" 3)"
	# assert the storage status
	assert_storage "alive" "$(life_status "data" 3)"
	# assert the storage label name
	assert_storage "data/3" "$(label 1)"
	# assert the unit attachment name
	assert_storage "dummy-storage-tp/0" "$(unit_attachment "data" 3 0)"
	# assert the attached unit state
	assert_storage "alive" "$(unit_state "data" 3 "dummy-storage-tp" 0)"
	echo "Filesystem tmpfs PASSED"

}

assess_fs() {
	echo "Assessing filesystem"
	juju deploy -m "${model_name}" ./testcharms/charms/dummy-storage-fs dummy-storage-np --series jammy --storage data=1G
	if [ $(unit_exist "data/4") ]; then
		return 0
	fi
	wait_for_storage "attached" '.storage["data/4"]["status"].current'
	# assert the storage kind name
	assert_storage "filesystem" "$(kind_name "data" 4)"
	# assert the storage status
	assert_storage "alive" "$(life_status "data" 4)"
	# assert the storage label name
	assert_storage "data/4" "$(label 2)"
	# assert the unit attachment name
	assert_storage "dummy-storage-np/0" "$(unit_attachment "data" 4 0)"
	# assert the attached unit state
	assert_storage "alive" "$(unit_state "data" 4 "dummy-storage-np" 0)"
	echo "Filesystem PASSED"
}

assess_multiple_fs() {
	echo "Assessing multiple filesystem, block, rootfs, loop"
	juju deploy -m "${model_name}" ./testcharms/charms/dummy-storage-fs dummy-storage-mp --series jammy --storage data=1G
	if [ $(unit_exist "data/5") ]; then
		return 0
	fi
	wait_for_storage "attached" '.storage["data/5"]["status"].current'
	# assert the storage kind name
	assert_storage "filesystem" "$(kind_name "data" 5)"
	# assert the storage status
	assert_storage "alive" "$(life_status "data" 5)"
	# assert the storage label name
	assert_storage "data/5" "$(label 3)"
	# assert the unit attachment name
	assert_storage "dummy-storage-mp/0" "$(unit_attachment "data" 5 0)"
	# assert the attached unit state
	assert_storage "alive" "$(unit_state "data" 5 "dummy-storage-mp" 0)"
	echo "Multiple filesystem, block, rootfs, loop PASSED"
}

# removes the applications if they exist.
remove_applications() {
	if [ $(application_exist "dummy-storage-fs") == "true" ]; then
		juju remove-application dummy-storage-fs --destroy-storage
	fi

	if [ $(application_exist "dummy-storage-lp") == "true" ]; then
		juju remove-application dummy-storage-lp --destroy-storage
	fi

	if [ $(application_exist "dummy-storage-tp") == "true" ]; then
		juju remove-application dummy-storage-tp --destroy-storage
	fi

	if [ $(application_exist "dummy-storage-np") == "true" ]; then
		juju remove-application dummy-storage-np --destroy-storage
	fi

	if [ $(application_exist "dummy-storage-mp") == "true" ]; then
		juju remove-application dummy-storage-mp --destroy-storage
	fi
}

# checks if the given storage unit exists.
unit_exist() {
	local name
	name=${1}
	juju storage --format json | jq "any(paths; .[-1] == \"${name}\")"
}

application_exist() {
	local name
	name=${1}
	juju status --format json | jq "any(paths; .[-1] == \"${name}\")"
}

run_charm_storage() {
	echo

	model_name="test-charm-storage"
	file="${TEST_DIR}/${model_name}.log"

	ensure "${model_name}" "${file}"

	echo "Assess create-storage-pool"
	juju create-storage-pool -m "${model_name}" loopy loop size=1G
	juju create-storage-pool -m "${model_name}" rooty rootfs size=1G
	juju create-storage-pool -m "${model_name}" tempy tmpfs size=1G
	juju create-storage-pool -m "${model_name}" ebsy ebs size=1G
	echo "create-storage-pool PASSED"

	# Assess the above created storage pools.
	echo "Assessing storage pool"
	if [ "${BOOTSTRAP_PROVIDER:-}" == "ec2" ]; then
		juju storage-pools -m "${model_name}" --format json | jq '.ebs | .provider' | check "ebs"
		juju storage-pools -m "${model_name}" --format json | jq '.["ebs-ssd"] | .provider' | check "ebs"
		juju storage-pools -m "${model_name}" --format json | jq '.tmpfs | .provider' | check "tmpfs"
		juju storage-pools -m "${model_name}" --format json | jq '.loop | .provider' | check "loop"
		juju storage-pools -m "${model_name}" --format json | jq '.rootfs | .provider' | check "rootfs"
	else
		juju storage-pools -m "${model_name}" --format json | jq '.rooty | .provider' | check "rootfs"
		juju storage-pools -m "${model_name}" --format json | jq '.tempy | .provider' | check "tmpfs"
		juju storage-pools -m "${model_name}" --format json | jq '.loopy | .provider' | check "loop"
		juju storage-pools -m "${model_name}" --format json | jq '.ebsy | .provider' | check "ebs"
	fi
	echo "Storage pool PASSED"

	assess_rootfs
	assess_loop_disk
	assess_tmpfs
	assess_multiple_fs

	remove_applications
	echo "All charm storage tests PASSED"

	destroy_model "${model_name}"
}

test_charm_storage() {
	if [ "$(skip 'test_charm_storage')" ]; then
		echo "==> TEST SKIPPED: charm storage tests"
		return
	fi

	(
		set_verbosity

		cd .. || exit

		run "run_charm_storage"
	)
}
