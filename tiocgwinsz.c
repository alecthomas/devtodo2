/* A hack to work around Go not having access to IOCTL constants. */
#include <stdio.h>
#include <sys/ioctl.h>

int main() {
  FILE *f = fopen("tiocgwinsz.go", "wt");
  fprintf(f, "package main\n");
  fprintf(f, "const TIOCGWINSZ = %lu\n", TIOCGWINSZ);
  fclose(f);
  return 0;
}
