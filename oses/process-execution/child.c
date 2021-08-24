#include <stdio.h>
#include <stdlib.h>
#include <sys/wait.h>
#include <string.h>
#include <unistd.h>

#define SLEEP_SEC 5

int main(int argc, char *argv[]) {
  int status;

  pid_t pid;
  pid = fork();

  if (pid == 0) {
    printf("Child PID: %d\nWaiting for a signal for %d seconds.\n", getpid(), SLEEP_SEC);
    sleep(SLEEP_SEC);
  }

  if ((pid = waitpid(pid, &status, 0)) > 0) {
    if (WIFEXITED(status)) {
      printf("Exited normally with status %d.\n", WEXITSTATUS(status));
    } else if (WIFSIGNALED(status)) {
      printf("Terminated via signal %d: %s.\n", WTERMSIG(status), strsignal(status));
    }
  }

  return 0;
}
