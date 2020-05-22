package com.fh;
public class MethodReloadTest {
    public static void main(String[] args) {
        int sum = add(100, 200);
        print(sum);
    }

    public static int add(int x, int y) {
        return x + y;
    }

    public static int add(int x, int y, int z) {
        return x + y + z;
    }

    public static native void print(int num);
}
