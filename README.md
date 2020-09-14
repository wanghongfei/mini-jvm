# Mini-JVM

使用Go实现的Java虚拟机(不完整)，解释执行，仅学习JVM使用；

Mini-JVM首先会从`classpath`中加载主类的class文件，然后找到main方法的字节码解释执行；执行过程中如果遇到新的类符号引用，则会通过全限定性名再从`classpath`中加载新的类文件，以此类推；

控制台输出、多线程功能通过自定义的标准库"mini-lib"中的`Printer`和`MiniThread`实现，没有使用JDK的标准库`Thread`，可以执行`compile-minilib.sh`编译mini-lib源文件；

当前支持的特性有：

- int加法
- 条件判断、for循环
- 控制台打印
- 简单对象(POJO)创建
- 基本类型数组和引用类型的数组创建、读写
- 字符串常量，即`String name = "hello, 世界"`
- main方法中可以读取到命令行参数
- 对象字段读写、静态字段读写
- 方法重载、方法重写、接口方法调用、形参全部为int类型的static方法调用
- native方法调用(本地方法表)
- 部分继承特性(字段继承、方法继承)
- 非标准库Thread类的线程支持
- synchronized关键字同步支持
- 支持部分Class方法，如toString(), getName(), isPrimitive()



  ![路线图](https://s1.ax1x.com/2020/06/30/NIpjTU.png)



## 编译运行Mini-JVM

编译(打开mod支持)：

```
go build -o mini-jvm
```

运行：

```shell
./mini-jvm -main [主类全限定性名，例如cn.fh.XXX] -classpath [类路径,可以是目录也可以是jar包路径, 多个用逗号分隔] [命令行参数,可选]
```

由于Mini-JVM的控制台输出和线程用的是私有类而JDK中`rt.jar`中的类，所以需要在classpath中指定`mini-lib`所在路径，例如：

```shell
./mini-jvm -main x.y.x.XXXMain -classpath mini-lib,classpath1,classpath2 cmd1 cmd2
```

单元测试`mini_jvm_test.go`中的case需要先修改`rtJarPath`为自己机器上`rt.jar`的路径后才能跑通：

```go
// 改成自己电脑中rt.jar的路径
var rtJarPath = "/Library/Java/JavaVirtualMachines/jdk1.8.0_181.jdk/Contents/Home/jre/lib/rt.jar"
```



## 编译testcase里的java代码

```shell
./compile-testcase.sh [.java文件相对于项目根目录的路径]
./compile-testcase.sh com/fh/ArrayTest.java
```





## 已支持的字节码

解释器已经支持的字节码如下：

```go
const (
	Nop byte = 0x00
	Aconstnull = 0x01

	Iconst0 = 0x03
	Iconst1 = 0x04
	Iconst2 = 0x05
	Iconst3 = 0x06
	Iconst4 = 0x07
	Iconst5 = 0x08

	Ldc = 0x12

	Iaload = 0x2e

	Aaload = 0x32
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

	Aload = 0x19
	Aload0 = 0x2a
	Aload1 = 0x2b
	Aload2 = 0x2c
	Aload3 = 0x2d

	Getstatic = 0xb2
	Putstatic = 0xb3

	Athrow = 0xbf

	Monitorenter = 0xc2
	Monitorexit = 0xc3

	Istore = 0x36
	Lstore1 = 0x40

	Astore = 0x3a
	Astore0 = 0x4b
	Astore1 = 0x4c
	Astore2 = 0x4d
	Astore3 = 0x4e
	Iastore = 0x4f

	Aastore = 0x53
	Castore = 0x55
	Pop = 0x57

	Dup = 0x59

	Iadd = 0x60
	Isub = 0x64

	Ishl = 0x78

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
	Ifacmpeq = 0xa5
	Ifacmpne = 0xa6
	Goto = 0xa7

	Areturn = 0xb0
	Return = 0xb1

	GetField = 0xb4
	Putfield = 0xb5

	Newarray = 0xbc
	Anewarray = 0xbd

	Invokevirtual = 0xb6
	Invokespecial = 0xb7
	Invokestatic = 0xb8
	Invokeinterface = 0xb9

	New = 0xbb

	Arraylength = 0xbe

	Ireturn = 0xac

	Wide = 0xc4
	Ifnonnull = 0xc7
)
```



## 已实现的特性举例

`mini_jvm_test.go`中有所有用例的单元测试；

以下Java代码均使用Java8进行编译；

  

著名的汉诺塔问题(`testcase/src/com/fh/Hanoi.java`)：

```java
package com.fh;
import cn.minijvm.io.Printer;

public class Hanoi {
    final static char A = 1;  //设置3个字符标记3根柱子，只是标记作用，实际计算中用不到。
    final static char B = 2;
    final static char C = 3;

    private int steps = 0;  //用于记录总共移动的次数

    public static void main(String[] args) {
        Hanoi hanoi = new Hanoi();
        hanoi.move(7, A, B, C);

        Printer.print(hanoi.getSteps());
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
}

```

```shell
./mini-jvm -main com.fh.Hanoi -classpath ../testcase/classes,../mini-lib/classes  # 输出127
```



"线程"支持：`testcase/src/com/fh/thread/ThreadTest.java`

```java
package com.fh.thread;

import cn.minijvm.concurrency.MiniThread;
import cn.minijvm.io.Printer;

public class ThreadTest {
    public static void main(String[] args) {
        Printer.print(args.length);
        for (String str : args) {
            Printer.printString(str);
        }

        // 创建协程支持的线程
        MiniThread th1 = new MiniThread();
        MiniThread th2 = new MiniThread();

        // 启动并执行线程
        MyTask task = new MyTask();
        th1.start(task);
        th2.start(task);

        // 当前线程休眠, 防止刚启动的线程还没来得及运行
        MiniThread.sleepCurrentThread(1);

        // 创建协程支持的线程
        MiniThread th3 = new MiniThread();
        MiniThread th4 = new MiniThread();
        YourTask task2 = new YourTask();
        th3.start(task2);
        th4.start(task2);

        // 当前线程休眠, 防止刚启动的线程还没来得及运行
        MiniThread.sleepCurrentThread(1);
    }

    public static class MyTask implements Runnable {
        private int number = 0;

        public void run() {
            for (int ix = 0; ix < 100; ix++) {
                synchronized (this) {
                    this.number++;
                    Printer.print(number);
                }
            }
        }
    }

    public static class YourTask implements Runnable {
        private int number = 0;

        public void run() {
            for (int ix = 0; ix < 100; ix++) {
                incr();
            }
        }

        private synchronized void incr() {
            this.number++;
            Printer.print(number);
        }
    }

}
```

输出：

```shell
=== RUN   TestThread
MiniJvm
1
2
3
...
...
200
--- PASS: TestThread (6.90s)
PASS
```



字符串打印：`testcase/src/com/fh/StringTest.java`

```java
package com.fh;
import cn.minijvm.io.Printer;

public class StringTest {
    public static void main(String[] args) {
        String name = "hello, 世界";
        Printer.printString(name);

        String dcs = new String("数字战斗模拟");
        Printer.printString(dcs);
    }
}

```

输出：

```shell
hello, 世界
数字战斗模拟
```





计算从1加到100的值(`testcase/src/com/fh/ForLoopPrintTest.java`)：

```java
package com.fh;

import cn.minijvm.io.Printer;

public class ForLoopPrintTest {
    public static void main(String[] args) {
        int sum = 0;
        for (int ix = 1; ix <= 100; ++ix) {
            sum = add(sum, ix);
        }

        Printer.print(sum);
    }

    public static int add(int x, int y) {
        return x + y;
    }
}

```





递归调用方法(`testcase/src/com/fh/RecursionTest`)：

```java
package com.fh;
import cn.minijvm.io.Printer;

public class RecursionTest {
    public static void main(String[] args) {
        // 从1打印到100
        foo(1);
    }

    public static void foo(int i) {
        if (i > 100) {
            return;
        }

        Printer.print(i);
        i++;
        foo(i);
    }
}
```



  

简单对象的创建/访问(`testcase/src/com/fh/NewSimpleObjectTest.java`)：

```java
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

```

   

方法重载(`testcase/src/com/fh/MethodReloadTest.java`)：

```java
package com.fh;
import cn.minijvm.io.Printer;

public class MethodReloadTest {
    public static void main(String[] args) {
        int sum = add(100, 200);
        Printer.print(sum);
    }

    public static int add(int x, int y) {
        return x + y;
    }

    public static int add(int x, int y, int z) {
        return x + y + z;
    }
}

```



  

方法重写(`testcase/src/com/ch/ClassExtendTest.java`)：

```java
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
        Printer.print(sum);
    }


```

  

int数组的读写： `testcase/src/com/fh/ArrayTest.java`



