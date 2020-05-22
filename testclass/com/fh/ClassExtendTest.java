package com.fh;

public class ClassExtendTest {
    public static void main(String[] args) {
        Person person = new Person();
        print(person.say());

        person = new Student();
        print(person.say());
    }

    public static native void print(int num);
}