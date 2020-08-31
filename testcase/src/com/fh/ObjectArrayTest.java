
package com.fh;

import cn.minijvm.io.Printer;

public class ObjectArrayTest {
    public static void main(String[] args) {
        Person[] people = new Person[2];
        people[0] = new Person();
        people[1] = new Person();

        Printer.print(people[0].getAge());
        people[0].setAge(100);
        Printer.print(people[0].getAge());
    }
}