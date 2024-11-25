using System;
using System.Runtime.InteropServices;

class Program
{
    [DllImport("fastLaunch.dll", CharSet = CharSet.Unicode, CallingConvention = CallingConvention.Cdecl)]
    public static extern void GoFunction(string input);

    static void Main(string[] args)
    {
        string input = "你好，世界";
        GoFunction(input);
    }
}
