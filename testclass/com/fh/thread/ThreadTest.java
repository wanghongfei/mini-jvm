package com.fh.thread;


public class ThreadTest {
    public static void main(String[] args) {
        executeInThread(new MyTask());

        print(-100);
        threadSleep(4);
        print(-200);
    }

    public static class MyTask implements Runnable {
        public void run() {
            for (int ix = 0; ix < 10; ix++) {
                print(ix);
            }
        }

        public native void print(int num);
    }

    /**
     * 在新的线程中执行task
     */
    public static native void executeInThread(Runnable task);

    /**
     * 休眠当前线程
     */
    public static native void threadSleep(int second);
    public static native void print(int num);
}
