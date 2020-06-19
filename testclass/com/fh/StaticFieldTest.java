package com.fh;

public class StaticFieldTest {
    private static int age = 100;
    private static int id = 200;

    public static void main(String[] args) {
        print(age);
        id = 400;
        print(id);
    }

    public static native void print(int num);
}