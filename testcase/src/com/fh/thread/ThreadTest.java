package com.fh.thread;

import cn.minijvm.concurrency.MiniThread;
import cn.minijvm.io.Printer;

public class ThreadTest {
    public static void main(String[] args) {
        // 创建协程支持的线程
        MiniThread th1 = new MiniThread();
        MiniThread th2 = new MiniThread();

        // 启动并执行线程
        MyTask task = new MyTask();
        th1.start(task);
        th2.start(task);

        // 当前线程休眠, 防止刚启动的线程还没来得及运行
        MiniThread.sleepCurrentThread(1);
    }

    public static class MyTask implements Runnable {
        private int number = 0;

        public void run() {
            for (int ix = 0; ix < 100; ix++) {
                synchronized (this) {
                    this.number++;
                    Printer.print(number);
                }
            }
        }
    }

}
