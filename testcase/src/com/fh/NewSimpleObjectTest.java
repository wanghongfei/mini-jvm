package com.fh;
import cn.minijvm.io.Printer;

public class NewSimpleObjectTest {
    public static void main(String[] args) {
        int sum = 0;
        for (int ix = 1; ix <= 100; ++ix) {
            sum = add(sum, ix);
        }

        Person p = new Person();
        p.setAge(sum);
        int age = p.getAge();

        Printer.print(age);
    }

    public static int add(int x, int y) {
        return x + y;
    }
}
