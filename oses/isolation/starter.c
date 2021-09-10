#define _GNU_SOURCE
#include <sched.h>
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <sys/wait.h>
#include <sys/prctl.h>
#include <sys/types.h>
#include <sys/capability.h>

#define MAX_MEM_BYTES   "2M"
#define MAX_PROCS       5
#define NUM_CAPS        5
#define NUM_CONTROLLERS 2
#define STACK_SIZE      65536

struct child_config {
  int argc;
  char **argv;
};

int capabilities() {
  int to_drop[NUM_CAPS] = {
    CAP_DAC_OVERRIDE,    // Disallow bypassing file permissions
    CAP_IPC_LOCK,        // Disallow memory locking
    CAP_SYS_BOOT,        // Disallow use of reboot, kexec_load
    CAP_SYS_ADMIN,       // Disallow sys admin capabilities
    CAP_SYS_MODULE       // Disallow loading/unloading kernel modules
  };

  // Remove capabilities from the bounding set
  for (int i = 0; i < NUM_CAPS; i++) {
    // for PR_CAPBSET_DROP, args 3 to 5 are unused
    if (prctl(PR_CAPBSET_DROP, to_drop[i], 0, 0, 0) == -1) {
      fprintf(stderr, "Error dropping capability %d", to_drop[i]);
      return -1;
    }
  }

  // Remove inheritable capabilities
  cap_t caps = cap_get_proc();
  if (caps == NULL) {
    fprintf(stderr, "Error getting caps with cap_get_proc()");
    return -1;
  }

  if (cap_set_flag(caps, CAP_INHERITABLE, NUM_CAPS, to_drop, CAP_CLEAR) == -1) {
    fprintf(stderr, "Error setting cap flags with cap_set_flag()");
    return -1;
  }

  if (cap_set_proc(caps) == -1) {
    fprintf(stderr, "Error setting inheritable caps with cap_set_proc()");
    return -1;
  }
  cap_free(caps);

  return 0;
}

int cgroup(pid_t pid) {
  // TODO: Make and cleanup directories in code instead of manually

  FILE *f;
  char path[1024];
  char *controllers[NUM_CONTROLLERS] = {"pids", "memory"};

  // Set up controllers by adding this task to the tasks file
  for (int i = 0; i < NUM_CONTROLLERS; i++) {
    sprintf(path, "/sys/fs/cgroup/%s/container/tasks", controllers[i]);
    f = fopen(path, "w");
    if (f == NULL) {
      fprintf(stderr, "Could not open %s", path);
      return -1;
    }
    fprintf(f, "%d", 0); // The value 0 adds the current process
    fclose(f);
  }

  // Set the max number of processes
  f = fopen("/sys/fs/cgroup/pids/container/pids.max", "w");
  if (f == NULL) {
    fprintf(stderr, "Could not open pids.max");
    return -1;
  }
  fprintf(f, "%d", MAX_PROCS);
  fclose(f);

  // Set the max memory
  f = fopen("/sys/fs/cgroup/memory/container/memory.limit_in_bytes", "w");
  if (f == NULL) {
    fprintf(stderr, "Could not open memory.limit_in_bytes");
    return -1;
  }
  fprintf(f, "%s", MAX_MEM_BYTES);
  fclose(f);

  return 0;
}

/* Entry point for child after `clone` */
int child(void *arg) {
  struct child_config *config = arg;

  // Set capabilities for child process
  if (capabilities() == -1) {
    fprintf(stderr, "Setting capabilities failed");
    return -1;
  }

  // Exec the program
  if (execvpe(config->argv[0], config->argv, NULL)) {
    fprintf(stderr, "execvpe failed %m.\n");
    return -1;
  }
  return 0;
}

int main(int argc, char**argv) {
  struct child_config config = {0};
  int flags = 0;
  pid_t child_pid = 0;

  // Prepare child configuration
  config.argc = argc - 1;
  config.argv = &argv[1];

  // Allocate stack for child
  char *stack = 0;
  if (!(stack = malloc(STACK_SIZE))) {
    fprintf(stderr, "Malloc failed");
    exit(1);
  }

  // Set up cgroup
  if (cgroup(getpid()) == -1) {
    fprintf(stderr, "Creating cgroup failed");
    return -1;
  }

  // Clone parent, enter child code
  if ((child_pid = clone(child, stack + STACK_SIZE, flags | SIGCHLD | CLONE_NEWNET | CLONE_NEWNS, &config)) == -1) {
    fprintf(stderr, "Clone failed");
    exit(2);
  }

  wait(NULL);

  return 0;
}
