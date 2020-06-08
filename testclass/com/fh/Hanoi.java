package com.fh;

public class Hanoi {
    final static char A = 1;  //设置3个字符标记3根柱子，只是标记作用，实际计算中用不到。
    final static char B = 2;
    final static char C = 3;

    private int steps = 0;  //用于记录总共移动的次数

    public static void main(String[] args) {
        Hanoi hanoi = new Hanoi();
        hanoi.move(7, A, B, C);

        printInt(hanoi.getSteps());
    }

    public void move(int n, char A, char B, char C) {
        if (n > 1) {
            move(n - 1, A, C, B);  //将 n - 1 个盘子从 A 移动到 B 上
            move(n - 1, B, C, A);  //将 n - 1 个盘子从 B 移动到 C 上
        }

        steps++;
    }

    public int getSteps() {
        return this.steps;
    }

    public static native void printInt(int num);
}

