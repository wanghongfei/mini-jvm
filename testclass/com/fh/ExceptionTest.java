package com.fh;

public class ExceptionTest {
    public static void main(String[] args) {
        try {
            print(10);
            throw new SimpleException();
        } catch (SimpleException e) {
            print(20);
        } finally {
            print(30);
        }
    }

    public static native void print(int num);
}