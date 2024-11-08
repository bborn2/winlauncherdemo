using System;
using System.Runtime.InteropServices;
 
class Program
{
    [DllImport("reminder.dll", CallingConvention = CallingConvention.Cdecl)]
    public static extern int Add(int a, int b);
 
    static void Main()
    {
        int result = Add(1,2);
        Console.WriteLine("Result: " + result);  // 输出结果应该是 8
        Console.ReadLine();
    }
}