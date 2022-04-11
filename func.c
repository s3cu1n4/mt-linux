// SPDX-License-Identifier: GPL-2.0

#include <stdio.h>
#include <string.h>
#include <unistd.h>
#include <signal.h>
#include <fcntl.h>

#include "func.h"


static volatile int g_rb_exit = 0;
static volatile int g_rb_quiet = 0;

#define MSG_SLOT_LEN    (65536)
static char g_rb_msg[MSG_SLOT_LEN] = {0};
static int rb_read_message(int fd)
{
    int rc, i, s = 0;
    rc = read(fd, g_rb_msg, MSG_SLOT_LEN - 1);
    for (i = 0; i <= rc; i++) {
        if (i == rc || g_rb_msg[i] == 0x17) {
            g_rb_msg[i] = 0;
            if (!g_rb_quiet && i > s) {
                PrintGo(&g_rb_msg[s]);
            }
                               
            s = i + 1;
        } else if (g_rb_msg[i] == 0x1e) {
            g_rb_msg[i] = ' ';
        }
    }
    memset(g_rb_msg, 0, MSG_SLOT_LEN);

    return rc;
}

static void rb_sigint_catch(int foo)
{
    printf("got CTRL + C, quiting ...\n");
    g_rb_exit = 1;
}

int Getsysmem()
{
    int fd = -1;
    signal(SIGINT, rb_sigint_catch);
    fd = open("/proc/elkeid-endpoint", O_RDONLY);
    if (fd < 0) {
        printf("Error: failed to open hids endpoint\n");
        return -1;
    }
    /* do consuming */
    while (!g_rb_exit) {
        rb_read_message(fd);
    }


    if (fd >= 0) {
        close(fd);
        printf("close(fd) ...\n");
        CheckMod();
    }
       
    return 0;
}
