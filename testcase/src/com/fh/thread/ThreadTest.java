package com.fh.thread;

import cn.minijvm.concurrency.MiniThread;
import cn.minijvm.io.Printer;

public class ThreadTest {
    public static void main(String[] args) {
        // 创建协程支持的线程
        MiniThread th = new MiniThread();
        Printer.print(-100);

        // 启动并执行线程
        th.start(new MyTask());
        // 当前线程休眠, 防止刚启动的线程还没来得及运行
        MiniThread.sleepCurrentThread(4);

        Printer.print(-200);
    }

    public static class MyTask implements Runnable {
        public void run() {
            for (int ix = 0; ix < 10; ix++) {
                Printer.print(ix);
            }
        }
    }

}
