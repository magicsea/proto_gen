using Microsoft.VisualStudio.TestTools.UnitTesting;
using msg_def;
using System;
using System.Collections.Generic;
using System.Text;
using Misc;
namespace msg.Tests
{
    [TestClass()]
    public class RoleTest
    {
        [TestMethod()]
        public void TestBs()
        {
            var bs = new BytesStream(new byte[128]);
            bs.WriteInt16(16);
            bs.WriteInt32(32);
            bs.WriteInt64(64);
            bs.WriteBool(true);
            bs.WriteByte(111);
            bs.WriteFloat(1.1f);
            bs.WriteDouble(2.2f);
            bs.WriteString("hello");
            bs.WriteUInt32(1234);
            bs.WriteBytes(new byte[]{ 1,2});
            bs.WriteUInt64(3721);


            Assert.AreEqual(16, bs.ReadInt16());
            Assert.AreEqual(32, bs.ReadInt32());
            Assert.AreEqual(64, bs.ReadInt64());
            Assert.AreEqual(true, bs.ReadBool());
            Assert.AreEqual(111, bs.ReadByte());
            Assert.AreEqual(1.1f, bs.ReadFloat());
            Assert.AreEqual(2.2f, bs.ReadDouble());
            Assert.AreEqual("hello", bs.ReadString());
            Assert.AreEqual((UInt32)(1234), bs.ReadUInt32());
            var bt = bs.ReadBytes();
            //Assert.AreEqual<byte[]>(bt, new byte[] { 1, 2 });
            Assert.AreEqual(Encoding.UTF8.GetString(bt), Encoding.UTF8.GetString(new byte[] { 1, 2 }));
            Assert.AreEqual<UInt64>(3721, bs.ReadUInt64());

            //var ret = Equals(new byte[] { 1, 2 }, new byte[] { 1, 2 });
            //Assert.IsTrue(ret);

            Assert.AreEqual(bs.RPos, bs.WPos);
        }


        [TestMethod()]
        public void TestVec()
        {
            var bs = new BytesStream(new byte[128]);
            var v = new MyVec2 { X = 1 ,Y = 2 };
            v.Write(bs);
            var v2 = new MyVec2();
            v2.Read(bs);
            Assert.AreEqual(v.X, v2.X);
            Assert.AreEqual(v.Y, v2.Y);
            Assert.AreEqual(bs.RPos, bs.WPos);
        }

        [TestMethod()]
        public void TestRole()
        {
            {//string
                var bs = new BytesStream(new byte[128]);
                var v = new MyVec2 { X = 1, Y = 2 };
                var role = new MyRole { Arri32 = new int[] { 3, 4 }, Bt = new byte[] { 5, 6 }, Myname = "haha", Mypos = v, Plist = new MyVec2[] { v }, SList = new string[] { "aaa", "bbb" }, sss = "hello" };

                role.Write(bs);

                var role2 = new MyRole();
                role2.Read(bs);
                Assert.AreEqual(bs.RPos, bs.WPos);
                Assert.AreEqual(role.Myname, role2.Myname);
                Assert.AreEqual(role.sss, role2.sss);

            }
            {//int
                var bs = new BytesStream(new byte[128]);
                var v = new MyVec2 { X = 1, Y = 2 };
                var role = new MyRole { Arri32 = new int[] { 3, 4 }, Bt = new byte[] { 5, 6 }, Myname = "haha", Mypos = v, Plist = new MyVec2[] { v }, SList = new string[] { "aaa", "bbb" }, sss = 1 };

                role.Write(bs);

                var role2 = new MyRole();
                role2.Read(bs);
                Assert.AreEqual(bs.RPos, bs.WPos);
                Assert.AreEqual(role.Myname, role2.Myname);
                Assert.AreEqual(role.sss, role2.sss);

            }
            {//class
                var bs = new BytesStream(new byte[128]);
                var v = new MyVec2 { X = 1, Y = 2 };
                var role = new MyRole { Arri32 = new int[] { 3, 4 }, Bt = new byte[] { 5, 6 }, Myname = "haha", Mypos = v, Plist = new MyVec2[] { v }, SList = new string[] { "aaa", "bbb" }, sss = new MyVec2 { X=1,Y=2 } };

                role.Write(bs);

                var role2 = new MyRole();
                role2.Read(bs);
                Assert.AreEqual(bs.RPos, bs.WPos);
                Assert.AreEqual(role.Myname, role2.Myname);
                var v1 = role.sss as MyVec2;
                var v2 = role2.sss as MyVec2;
                Assert.AreEqual(v1.X,v2.X);

            }
            {//byte[]
                var bs = new BytesStream(new byte[128]);
                var v = new MyVec2 { X = 1, Y = 2 };
                var role = new MyRole { Arri32 = new int[] { 3, 4 }, Bt = new byte[] { 5, 6 }, Myname = "haha", Mypos = v, Plist = new MyVec2[] { v }, SList = new string[] { "aaa", "bbb" }, sss = new byte[] { 1,2} };

                role.Write(bs);

                var role2 = new MyRole();
                role2.Read(bs);
                Assert.AreEqual(bs.RPos, bs.WPos);
                Assert.AreEqual(role.Myname, role2.Myname);
                var v1 = (byte[])role.sss;
                var v2 = (byte[])role2.sss;
                Assert.AreEqual(v1[0], v2[0]);
                Assert.AreEqual(v1[1], v2[1]);

            }
        }

        [TestMethod()]
        public void TestPackArgs()
        {
            var v = new MyVec2 { X = 3, Y = 2 };
            var list = new object[]{ 1, "sss", v};
            var data = MessagePacker.PackArgs(list);
            var ret = MessagePacker.UnPackArgs(data);
            for (int i = 0; i < 1; i++)
            {
                Assert.AreEqual(list[i], ret[i]);
            }
            Assert.AreEqual(v.X, (ret[2] as MyVec2).X);
        }


        [TestMethod()]
        public void TestGrowBuff()
        {
            {//<right
                var bs = new BytesStream(new byte[128]);
                bs.SetWriteGrow(true);
                bs.WPos = 127;
                bs.WriteInt32(888);
                Assert.AreEqual(256, bs.Capicity);
                bs.RPos = 127;
                Assert.AreEqual<Int32>(888, bs.ReadInt32());
            }

            {//<MAX_STREAM_GROW_MIN_SIZE
                var bs = new BytesStream(new byte[4]);
                bs.SetWriteGrow(true);
                bs.WriteInt32(32);
                Assert.AreEqual(4, bs.Capicity);
                bs.WriteInt32(1);
                Assert.AreEqual(BytesStream.STREAM_GROW_MIN_SIZE + 4, bs.Capicity);
            }

            {//needSize > STREAM_GROW_MIN_SIZE*2
                var bs = new BytesStream(new byte[128]);
                bs.SetWriteGrow(true);
                bs.WPos = 127;
                bs.WriteBytes(new byte[300]);
                Assert.AreEqual(559, bs.Capicity);
            }

            {//needSize > MAX_STREAM_SIZE/2
                var bs = new BytesStream(new byte[128]);
                bs.SetWriteGrow(true);
                bs.WPos = 127;
                bs.WriteBytes(new byte[BytesStream.MAX_STREAM_SIZE / 2]);
                Assert.AreEqual(8388608, bs.Capicity);
            }

        }
    }



}