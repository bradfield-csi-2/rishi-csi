/** @file
 *
 * DictFS Example
 *
 * Compile with:
 *
 *   gcc -Wall dictfs.c `pkg-config fuse3 --cflags --libs` -o dictfs
 *
 */

#define FUSE_USE_VERSION 31
#define MAX_FILES        10

#include <errno.h>
#include <fuse.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

struct node {
  char *name;
  struct node* nodes[MAX_FILES];
};

struct node* newNode(char *name)
{
  struct node* node
    = (struct node*)malloc(sizeof(struct node));

  node->name = name;
  return node;
}

struct node *dict;

static void *dict_init(struct fuse_conn_info *conn, struct fuse_config *cfg) {
 struct node *root = newNode("");
 root->nodes[0] = newNode("e4");
 root->nodes[1] = newNode("d4");

 dict = root;
 return NULL;
}

static int dict_getattr(const char *path, struct stat *stbuf, struct fuse_file_info *fi) {
  int res = 0;

  printf("%s\n", dict->name);

  stbuf->st_mode = S_IFDIR | 0755;
  stbuf->st_nlink = 2;

  return res;
}

static int dict_readdir(const char *path, void *buf, fuse_fill_dir_t filler, off_t offset, struct fuse_file_info *fi, enum fuse_readdir_flags flags) {
  filler(buf, ".", NULL, 0, 0);
  filler(buf, "..", NULL, 0, 0);

  // Walk the directory path
  struct node *cursor = dict;
  char *delim = "/";
  char *path_str = strdup(path);
  char *token = strtok(path_str, delim);
  int found;
  while (token) {
    found = 0;
    for (int i = 0; i < MAX_FILES; i++) {
      struct node *currNode = cursor->nodes[i];

      // Found a child element that matches that path segment, traverse into it
      // and break out of the for loop.
      if (currNode != NULL && strcmp(token, currNode->name) == 0) {
        cursor = currNode;
        token = strtok(NULL, delim);
        found = 1;
        break;
      }
    }
  }
  free(path_str);
  if (found) {
    return -ENOENT;
  }

  // Fill up the directory listing with the nodes at this level
  for (int i = 0; i < MAX_FILES; i++) {
    struct node *currNode = cursor->nodes[i];
    if (currNode != NULL) {
      filler(buf, cursor->nodes[i]->name, NULL, 0, 0);
    }
  }

  return 0;
}

static const struct fuse_operations dict_oper = {
  .init    = dict_init,
  .getattr = dict_getattr,
  .readdir = dict_readdir,
};

int main(int argc, char *argv[]) {
  int ret;
  struct fuse_args args = FUSE_ARGS_INIT(argc, argv);

  ret = fuse_main(args.argc, args.argv, &dict_oper, NULL);
  fuse_opt_free_args(&args);
  return ret;
}
