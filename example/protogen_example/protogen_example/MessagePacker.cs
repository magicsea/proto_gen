using System;
using System.Collections.Generic;
using System.Text;
using System.Reflection;

namespace Misc
{
    public enum ArgType : byte
    {
        Int8 = 1,
        UInt8 = 2,
        Int16 = 3,
        UInt16 = 4,
        Int32 = 5,
        UInt32 = 6,
        Int64 = 7,
        UInt64 = 8,
        Float32 = 9,
        Float64 = 10,
        String = 11,
        ByteArray = 12,
        Bool = 13,
        Nil = 14,
        Error = 15,

        StreamMsg = 20,
        ProtoBuffMsg = 21,
    }
    public struct Error
    {
        public string Message;
        public override string ToString()
        {
            return Message;
        }
    }
    public class MessagePacker
    {
        public static byte[] marshalProtobuf(object v)
        {
            var o = v as Google.Protobuf.IMessage;
            var data = Google.Protobuf.MessageExtensions.ToByteArray(o);
            return data;
        }

        public static object unMarshalProtobuf(Type t, byte[] data)
        {
            FieldInfo parserField = t.GetField("_parser", BindingFlags.Static | BindingFlags.NonPublic);
            object objParser = parserField.GetValue(null);
            MethodInfo method = parserField.FieldType.GetMethod("ParseFrom", new Type[] { typeof(byte[]) });
            object o = method.Invoke(objParser, new object[] { data });
            return o;
        }


        static byte[] tempBuf = new byte[1024 * 1024];
        static BytesStream tempStream = new BytesStream(tempBuf);
        public static byte[] PackArgs(object[] args)
        {
            tempStream.Reset();
            foreach (var item in args)
            {
                PackAny(item, tempStream);
            }

            var len = tempStream.WPos;

            var ret = new byte[len];
            Array.Copy(tempBuf, 0, ret, 0, len);
            return tempBuf;
        }

        public static object[] UnPackArgs(byte[] data)
        {
            List<object> rets = new List<object>();
            var readStream = new BytesStream(data);
            for (; ; )
            {
                try
                {
                    var v = UnPackAny(readStream);
                    rets.Add(v);
                    if (readStream.IsReadEnd())
                    {
                        break;
                    }
                }
                catch
                {
                    break;
                }
            }
            return rets.ToArray();
        }

        public static void PackAny(object o,BytesStream s)
        {
            if (o==null)
            {
                s.WriteByte((byte)ArgType.Nil);
                return;
            }

            var t = o.GetType();

            if (t.Equals(typeof(byte)))
            {
                s.WriteByte((byte)ArgType.Int8);
                s.WriteByte((byte)o);
            }
            else if (t.Equals(typeof(Int16)))
            {
                s.WriteByte((byte)ArgType.Int16);
                s.WriteInt16((Int16)o);
            }
            else if (t.Equals(typeof(UInt16)))
            {
                s.WriteByte((byte)ArgType.UInt16);
                s.WriteUInt16((UInt16)o);
            }
            else if (t.Equals(typeof(Int32)))
            {
                s.WriteByte((byte)ArgType.Int32);
                s.WriteInt32((Int32)o);
            }
            else if (t.Equals(typeof(UInt32)))
            {
                s.WriteByte((byte)ArgType.UInt32);
                s.WriteUInt32((UInt32)o);
            }
            else if (t.Equals(typeof(Int64)))
            {
                s.WriteByte((byte)ArgType.Int64);
                s.WriteInt64((Int64)o);
            }
            else if (t.Equals(typeof(UInt64)))
            {
                s.WriteByte((byte)ArgType.UInt64);
                s.WriteUInt64((UInt64)o);
            }
            else if (t.Equals(typeof(float)))
            {
                s.WriteByte((byte)ArgType.Float32);
                s.WriteFloat((float)o);
            }
            else if (t.Equals(typeof(double)))
            {
                s.WriteByte((byte)ArgType.Float64);
                s.WriteDouble((double)o);
            }
            else if (t.Equals(typeof(string)))
            {
                s.WriteByte((byte)ArgType.String);
                s.WriteString((string)o);
            }
            else if (t.Equals(typeof(byte[])))
            {
                s.WriteByte((byte)ArgType.ByteArray);
                s.WriteBytes((byte[])o);
            }
            else if (t.Equals(typeof(bool)))
            {
                s.WriteByte((byte)ArgType.Bool);
                s.WriteBool((bool)o);
            }
            else if (t.Equals(typeof(Error)))
            {
                s.WriteByte((byte)ArgType.Error);
                var err = (Error)o;
                s.WriteString(err.Message);
            }
            else
            {
                var pm = o as IPBMessage;
                if(pm!=null)
                {
                    var data = Google.Protobuf.MessageExtensions.ToByteArray(pm);
                    s.WriteByte((byte)ArgType.ProtoBuffMsg);
                    s.WriteUInt16(pm.MsgID);
                    s.WriteBytes(data);
                    return;
                }

                var sm = o as IMessage;
                if (sm != null)
                {
                    s.WriteByte((byte)ArgType.StreamMsg);
                    s.WriteUInt16(sm.MsgID);
                    sm.Write(s);
                    return;
                }

                throw new Exception("not support type " + t.Name);
            }
        }

        public static object UnPackAny(BytesStream s)
        {
            var typ = (ArgType)s.ReadByte();
            if (typ==0)
            {
                return null;
            }
            switch (typ)
            {
                case ArgType.Int8:
                case ArgType.UInt8:
                    {
                        return s.ReadByte();
                    }
                case ArgType.Int16:
                    {
                        return  s.ReadInt16();
                    }
                case ArgType.UInt16:
                    {
                        return s.ReadUInt16();
                    }
                case ArgType.Int32:
                    {
                        return s.ReadInt32();
                    }
                case ArgType.UInt32:
                    {
                        return s.ReadUInt32();
                    }
                case ArgType.Int64:
                    {
                        return s.ReadInt64();
                    }
                case ArgType.UInt64:
                    {
                        return s.ReadUInt64();
                    }
                case ArgType.String:
                    {
                        return s.ReadString();
                    }
                case ArgType.ByteArray:
                    {
                        return s.ReadBytes();
                    }
                case ArgType.Bool:
                    {
                        return s.ReadBool();
                    }
                case ArgType.Float32:
                    {
                        return s.ReadFloat();
                    }
                case ArgType.Float64:
                    {
                        return s.ReadDouble();
                    }
                case ArgType.Error:
                    {
                        return new Error { Message = s.ReadString()};
                    }
                case ArgType.ProtoBuffMsg:
                    {
                        var pid = s.ReadUInt16();
                        //Type pt = ProtoDef.Inst.GetTypeByID(pid);
                        //var pd = s.ReadBytes();
                        //if (pt != null)
                        //{
                        //    return unMarshalProtobuf(pt, pd);
                        //}
                        //else
                        //{
                        //    return new RawProto(pid, pd);
                        //}
                        throw new Exception("todo..pb");

                    }
                case ArgType.StreamMsg:
                    {
                        var pid = s.ReadUInt16();
                        var sm = msg_def.msg_def.Newmsg_def(pid);
                        sm.Read(s);
                        return sm;
                    }
            }

            throw new Exception("not support type:" + typ);
        }


        //todo:是不是要用ref
        public static int GetSizeAny(object msg)
        {
            //include type: 1 byte
            return 1 + GetDataSize(msg);
        }

        public static int GetDataSize(object msg)
        {
            bool foundBase = false;
            var ds = BytesStream.GetStreamBaseTypeSize(msg,out foundBase);
            if (foundBase)
            {
                return ds;
            }

            if (msg==null)
            {
                return 0;
            }

            var t = msg.GetType();
            if (t.Equals(typeof(Error)))
            {
                var err = (Error)msg;
                bool tmp;
                var errLen = BytesStream.GetStreamBaseTypeSize(err.Message, out tmp);
                return errLen;
            }
            var pm = msg as IPBMessage;
            if (pm!=null)
            {
                return 2 + pm.CalculateSize();
            }

            var sm = msg as IMessage;
            if (sm != null)
            {
                return 2 + sm.GetStreamSize();
            }
            throw new Exception("not support type:" + t);
        }
    }
}
