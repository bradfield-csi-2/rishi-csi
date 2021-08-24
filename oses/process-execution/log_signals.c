#include <signal.h>
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <string.h>

static void handler(int signum) {
  printf("Received Signal %d: %s!\n", signum, strsignal(signum));
}

int main() {
  struct sigaction sa;

  sa.sa_handler = handler;
  sigemptyset(&sa.sa_mask);
  sa.sa_flags = SA_RESTART; // Restart functions if interrupted by handler

  // Set handler for all signal numbers between 1 and SIGSYS (31)
  for (int i = 1; i <= SIGSYS; i++) {
    // These signals can be neither caught nor ignored
    if (i == SIGKILL || i == SIGSTOP) continue;

    if (sigaction(i, &sa, NULL) == -1) {
      printf("Error setting signal: %d\n", i);
      exit(1);
    }
  }

  printf("PID: %d\nWaiting for a signal...\n", getpid());
  while(1);
}
