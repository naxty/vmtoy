# vmtoy

vmtoy is a toy project to lightweight VMs on OSX. It is a command line tool that can be used to create, start, stop, and delete VMs. It will be written in Golang and use QEMU as virtualization software.


The tool will provide a CLI to:
- Create a VM
- Start a VM
- Stop a VM
- Delete a VM
- List all VMs
- Show VM details
- Show VM status
- Show VM logs
- SSH to VM
- SCP to VM
- Run a command on VM


## Idea

Providing golang implementation on top of QEMU to manage VM. Firstly, focus on Alpine Linux VMs. Later, extend to other Linux distros.

Alpine Linux implementation:

1. Check if the ISO is downloaded
2. Creating a image file `qemu-img create -f qcow2 alpine.qcow2 8G`
3. Running ISO and setting up the alpine linux:
   ```bash
   qemu-system-x86_64 \
    -name alpine-vm \
    -m 1024 \
    -smp 2 \
    -hda alpine.qcow2 \
    -cdrom /Users/stheno/git/dev/gofun/containers/alpine-virt-3.21.0-x86_64.iso \
    -boot d \
    -accel tcg \
    -serial telnet:localhost:4321,server,nowait \
    -nographic \
    -netdev user,id=user.0,hostfwd=tcp::2222-:22 \
    -device e1000,netdev=user.0
   ```
    - Trying to use telnet to connect to the VM
    - Setting up image via telnet commands
    - Poweroff the VM
4. Poweron again with the image file `qemu-system-x86_64 -name alpine-vm -m 1024 -smp 2 -hda alpine.qcow2 -accel tcg -nographic -netdev user,id=user.0,hostfwd=tcp::2222-:22 -device e1000,netdev=user.0 -serial /dev/pts/1`
5. SSH to VM `ssh -p 2222 root@localhost`



## Dependencies

```bash
brew install dosfstools
brew install cdrtools
```