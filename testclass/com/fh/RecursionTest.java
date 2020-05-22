package com.fh;

public class RecursionTest {
    public static void main(String[] args) {
        // 从1打印到100
        foo(1);
    }

    public static void foo(int i) {
        if (i > 100) {
            return;
        }

        print(i);
        i++;
        foo(i);
    }

    public static native void print(int num);
}