package com.fh;

import cn.minijvm.io.Printer;

public class ClassExtendTest {
    public static void main(String[] args) {
        Person person = new Person();
        Printer.print(person.say());

        person = new Student(); // Student继承了Person
        Printer.print(person.say());
    }
}