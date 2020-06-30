package cn.minijvm.concurrency;

public class MiniThread {
    public native void start(Runnable task);
    public static native void sleepCurrentThread(int second);
}