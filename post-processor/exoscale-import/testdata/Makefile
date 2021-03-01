# This target creates a dummy template disk file in QCOW2 format
# used to mock an artifact from the `qemu` Packer builder.
all:
	qemu-img create test-packer-post-processor-exoscale-import.qcow2.tmp 10240M
	qemu-img convert -c -O qcow2 \
		test-packer-post-processor-exoscale-import.qcow2.tmp \
		test-packer-post-processor-exoscale-import.qcow2
	rm -f test-packer-post-processor-exoscale-import.qcow2.tmp
