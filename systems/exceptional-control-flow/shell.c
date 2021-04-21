#include <errno.h>
#include <signal.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#define MAXLINE  1024
#define MAXARGS  28
#define EXIT_MSG "Bye for now!"

typedef void (*sighandler_t)(int);

// Function Prototypes
void eval(char *cmdline);
void parseline(char *buf, char *argv[]);
int builtin_cmd(char *argv[]);

// Signal Handlers
void sigint_handler(int sig) {
  return;
}

int main(int argc, char *argv[]) {
  char cmdline[MAXLINE];

  // Install signal handlers
  if (signal(SIGINT, sigint_handler) == SIG_ERR) {
    printf("signal error\n");
  }

  while (1) {
    printf("🦊 ");
    fgets(cmdline, MAXLINE, stdin);
    if (feof(stdin)) {
      printf("\n%s\n", EXIT_MSG);
      exit(0);
    }

    eval(cmdline);
  }
  return 0;
}

void eval(char *cmdline) {
  char *argv[MAXARGS];
  char buf[MAXLINE];
  pid_t pid;

  strcpy(buf, cmdline);
  parseline(buf, argv);
  if (argv[0] == NULL) {
    return;
  }

  if (builtin_cmd(argv)) {
    return;
  }

  if ((pid = fork()) == 0) {
    if (execvp(argv[0], argv) < 0) {
      printf("%s: Command not found.\n", argv[0]);
      exit(0);
    }
  }

  int status;
  if (waitpid(pid, &status, 0) < 0) {
    printf("waitfg: waitpid error\n");
  }
  return;
}

int builtin_cmd(char *argv[]) {
  if (!strcmp(argv[0], "exit")) {
    printf("%s\n", EXIT_MSG);
    exit(0);
  } else if (!strcmp(argv[0], "cd")) {
    // Change to home directory on empty cd
    char *dir = argv[1] == NULL ? getenv("HOME") : argv[1];
    if (chdir(dir) < 0) {
      printf("%s: %s\n", argv[0], strerror(errno));
      return 1;
    }
  }
  return 0;
}

// Taken from CS:APP, 3e, Figure 8.25
void parseline(char *buf, char *argv[]) {
  char *delim;
  int argc;

  buf[strlen(buf)-1] = ' '; // Replace trailing '\n' with space

  argc = 0;
  while ((delim = strchr(buf, ' '))) {
    argv[argc++] = buf;
    *delim = '\0'; // Replace the space with a null byte
    buf = delim + 1;
    // Skip extra spaces
    while (*buf && (*buf == ' ')) {
      buf++;
    }
  }
  argv[argc] = NULL;  // Need to have a NULL after all the arguments
  return;
}
