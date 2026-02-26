sudo virt-install --connect qemu:///system \
  --import \
  --name control01 \
  --ram 2048 --vcpus 2 \
  --os-variant=unknown \
  --disk size=20,backing_store=/var/lib/libvirt/images/flatcar-linux/flatcar_production_qemu_image.img \
  --graphics=none \
  --qemu-commandline="-fw_cfg name=opt/org.flatcar-linux/config,file=$(pwd)/control01.ign"