#include <stdio.h>
#include <stdlib.h>
#include <string.h>

int main(int argc, char *argv[]) {
	printf("Mini Obj Dump\n");
  FILE *fileptr;
  char *buf;
  long filelen;

  fileptr = fopen("/bin/true", "rb");   // Open the file in binary mode
  fseek(fileptr, 0, SEEK_END);          // Jump to the end of the file
  filelen = ftell(fileptr);             // Get the current byte offset in the file
  rewind(fileptr);                      // Jump back to the beginning of the file

  buf = (char *)malloc(filelen * sizeof(char)); // Enough memory for the file
  fread(buf, filelen, 1, fileptr); // Read in the entire file
  fclose(fileptr); // Close the file

  int magic;  // 4 byte magic number
  long offset;
  memcpy(&magic, buf, sizeof(magic));
  printf("Magic Number: %#08x\n", magic);
  offset = sizeof(magic);
  fseek(fileptr, offset, SEEK_SET);



	return 0;
}
