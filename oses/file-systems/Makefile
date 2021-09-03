build:
	gcc -Wall dictfs.c `pkg-config fuse3 --cflags --libs` -o dictfs

mount:
	./dictfs test

unmount:
	fusermount -u test
