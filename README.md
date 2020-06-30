# Mini-JVM

使用Go实现的Java虚拟机(不完整)，解释执行，仅学习JVM使用；

Mini-JVM首先会从`classpath`中加载主类的class文件，然后找到main方法的字节码解释执行；执行过程中如果遇到新的类符号引用，则会通过全限定性名再从`classpath`中加载新的类文件，以此类推；

控制台输出、多线程功能通过自定义的标准库"mini-lib"中的`Printer`和`MiniThread`实现，没有使用JDK的标准库`Thread`，可以执行`compile-minilib.sh`编译mini-lib源文件；

当前支持的特性有：

- int加法
- 循环结构
- 控制台输出
- 简单对象创建
- 对象字段读写、静态字段读写
- 方法重载、方法重写、接口方法调用、形参全部为int类型的static方法调用
- native方法调用(本地方法表)
- 部分继承特性(字段继承、方法继承)
- 非标准库Thread类的线程支持



  ![路线图](https://s1.ax1x.com/2020/06/11/tHM4OJ.png)



## 使用方法

编译(打开mod支持)：

```
go build -o mini-jvm
```

运行：

```shell
./mini-jvm [主类全限定性名] [classpath(支持.class文件所在目录或jar所在目录)]
```

由于Mini-JVM的控制台输出和线程用的是私有类而JDK中`rt.jar`中的类，所以需要在classpath中指定`mini-lib`所在路径，例如：

```shell
./mini-jvm x.y.x.XXXMain mini-lib/  classpath1/ classpath2/ ... ...
```





## 已支持的字节码

解释器已经支持的字节码如下：

```go
const (
	Nop byte = 0x00

	Iconst0 = 0x03
	Iconst1 = 0x04
	Iconst2 = 0x05
	Iconst3 = 0x06
	Iconst4 = 0x07
	Iconst5 = 0x08

	Iaload = 0x2e

	Caload = 0x34

	Istore0 = 0x3b
	Istore1 = 0x3c
	Istore2 = 0x3d
	Istore3 = 0x3e

	Bipush = 0x10
	Sipush = 0x11

	Iload = 0x15
	Iload0 = 0x1a
	Iload1 = 0x1b
	Iload2 = 0x1c
	Iload3 = 0x1d

	Aload0 = 0x2a
	Aload1 = 0x2b
	Aload2 = 0x2c

	Getstatic = 0xb2
	Putstatic = 0xb3

	Athrow = 0xbf

	Istore = 0x36
	Lstore1 = 0x40

	Astore0 = 0x4b
	Astore1 = 0x4c
	Astore2 = 0x4d
	Iastore = 0x4f

	Castore = 0x55

	Dup = 0x59

	Iadd = 0x60
	Isub = 0x64

	Iinc = 0x84

	Ifeq = 0x99
	Ifne = 0x9a
	Iflt = 0x9b
	Ifge = 0x9c
	Ifgt = 0x9d
	Ifle = 0x9e

	Ificmpeq = 0x9f
	Ificmpne = 0xa0
	Ificmplt = 0xa1
	Ificmpge = 0xa2
	Ificmpgt = 0xa3
	Ificmple = 0xa4
	Goto = 0xa7

	Return = 0xb1

	GetField = 0xb4
	Putfield = 0xb5

	Newarray = 0xbc

	Invokevirtual = 0xb6
	Invokespecial = 0xb7
	Invokestatic = 0xb8
	Invokeinterface = 0xb9

	New = 0xbb

	Ireturn = 0xac

	Wide = 0xc4
)
```



## 已实现的特性举例

`mini_jvm_test.go`中有所有用例的单元测试；

以下Java代码均使用Java8进行编译；

  

著名的汉诺塔问题(`testcase/src/com/fh/Hanoi.java`)：

```java
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

```

```shell
./mini-jvm cn.fh.Hanoi(主类名) testcase/classes/ mini-lib/classes  # 输出127
```



"线程"支持：`testcase/src/com/fh/thread/ThreadTest.java`

```java
package com.fh.thread;

import cn.minijvm.concurrency.MiniThread;
import cn.minijvm.io.Printer;

public class ThreadTest {
    public static void main(String[] args) {
        // 创建协程支持的线程
        MiniThread th = new MiniThread();
        Printer.print(-100);

        // 启动并执行线程
        th.start(new MyTask());
        // 当前线程休眠, 防止刚启动的线程还没来得及运行
        MiniThread.sleepCurrentThread(4);

        Printer.print(-200);
    }

    public static class MyTask implements Runnable {
        public void run() {
            for (int ix = 0; ix < 10; ix++) {
                Printer.print(ix);
            }
        }
    }

}

```

输出：

```shell
=== RUN   TestThread
-100
0
1
2
3
4
5
6
7
8
9
-200
--- PASS: TestThread (6.90s)
PASS
```



计算从1加到100的值(`testcase/src/com/fh/ForLoopPrintTest.java`)：

```java
package com.fh;
public class ForLoopPrintTest {
    public static void main(String[] args) {
        int sum = 0;
        for (int ix = 1; ix <= 100; ++ix) {
            sum = add(sum, ix);
        }
        print(sum);
    }

    public static int add(int x, int y) {
        return x + y;
    }

    public static native void print(int num);
}

```





递归调用方法(`testcase/src/com/fh/RecursionTest`)：

```java
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
```



  

简单对象的创建/访问(`testcase/src/com/fh/NewSimpleObjectTest.java`)：

```java
package com.fh;
public class NewSimpleObjectTest {
    public static void main(String[] args) {
        int sum = 0;
        for (int ix = 1; ix <= 100; ++ix) {
            sum = add(sum, ix);
        }

        Person p = new Person();
        p.setAge(sum);
        int age = p.getAge();

        print(age);
    }

    public static int add(int x, int y) {
        return x + y;
    }

    public static native void print(int num);
}

```

   

方法重载(`testcase/src/com/fh/MethodReloadTest.java`)：

```java
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

```



  

方法重写(`testcase/src/com/ch/ClassExtendTest.java`)：

```java
package com.fh;

public class ClassExtendTest {
    public static void main(String[] args) {
        Person person = new Person();
        print(person.say());

        person = new Student(); // Student继承了Person
        print(person.say());
    }

    public static native void print(int num);
}
```





判断语句、自增语句：

```java
package com.fh;

public class IfTest {
    public static void main(String[] args) {
        int sum = 0;

        if (sum > 0) {
            sum = 100;
        }
        if (sum >= 0) {
            sum -= 200;
        }
        if (sum < 0) {
            sum -= 100;
        }
        if (sum <= 0) {
            sum -= 1;
        }
        if (sum == 0) {
            sum += 2000;
        }

        print(sum); // -301

        if (sum > 10) {

        }
        if (sum >= 10) {

        }
        if (sum < 10) {

        }
        if (sum <= 10) {

        }
        if (sum == 10) {

        }
        print(sum);
    }

    public static native void print(int num);

```

  

int数组的读写： `testcase/src/com/fh/ArrayTest.java`



