#include <signal.h>
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <sys/wait.h>
#include <sys/socket.h>
#include <sys/un.h>
#include <unistd.h>

int START = 2, END = 20;
char *TESTS[] = {"brute_force", "brutish", "miller_rabin"};
int num_tests = sizeof(TESTS) / sizeof(char *);

int sock() {
  unsigned int s, s2;
  struct sockaddr_un local, remote;
  int len;

  s = socket(AF_UNIX, SOCK_STREAM, 0);

  local.sun_family = AF_UNIX;
  strcpy(local.sun_path, "/home/ubuntu/mysocket");
  unlink(local.sun_path);
  len = strlen(local.sun_path) + sizeof(local.sun_family);
  bind(s, (struct sockaddr *)&local, len);

  listen(s, 5);

  for(;;) {
    int done, n;
    printf("Waiting for a connection...\n");
    t = sizeof(remote);
    if ((s2 = accept(s, (struct sockaddr *)&remote, &t)) == -1) {
      perror("accept");
      exit(1);
    }

    printf("Connected.\n");

    done = 0;
    do {
      n = recv(s2, str, 100, 0);
      if (n <= 0) {
        if (n < 0) perror("recv");
        done = 1;
      }

      if (!done)
        if (send(s2, str, n, 0) < 0) {
          perror("send");
          done = 1;
        }
    } while (!done);

    close(s2);
  }
}

int main(int argc, char *argv[]) {
  int testfds[num_tests][2];
  int resultfds[num_tests][2];
  int result, i;
  long n;
  pid_t pid;

  for (i = 0; i < num_tests; i++) {
    pipe(testfds[i]);
    pipe(resultfds[i]);

    pid = fork();

    if (pid == -1) {
      fprintf(stderr, "Failed to fork\n");
      exit(-1);
    }

    if (pid == 0) {
      // we are the child, connect the pipes correctly and exec!
      close(testfds[i][1]);
      close(resultfds[i][0]);
      dup2(testfds[i][0], STDIN_FILENO);
      dup2(resultfds[i][1], STDOUT_FILENO);
      execl("primality", "primality", TESTS[i], (char *)NULL);
    }

    // we are the parent
    close(testfds[i][0]);
    close(resultfds[i][1]);
  }

  // for each number, run each test
  for (n = START; n <= END; n++) {
    for (i = 0; i < num_tests; i++) {

      // we are the parent, so send test case to child and read results
      write(testfds[i][1], &n, sizeof(n));
      read(resultfds[i][0], &result, sizeof(result));
      printf("%15s says %ld %s prime\n", TESTS[i], n, result ? "is" : "IS NOT");
    }
  }
}
