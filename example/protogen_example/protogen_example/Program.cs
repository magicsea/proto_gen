using System;
using msg_def;
using Misc;
namespace protogen_example
{
    partial class Program
    {
        static void Main(string[] args)
        {
            var bs = new BytesStream(new byte[128]);
            var v = new MyVec2 { X = 1, Y = 2 };
            v.Write(bs);
            var v2 = new MyVec2();
            v2.Read(bs);

            Console.WriteLine("Hello World!"+v2.X+v2.Y);
        }
    }

}


