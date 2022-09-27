using System;
using System.Collections.Generic;
using System.Text;
using System.IO;

namespace Misc
{
    public class BytesStream
    {
        public byte[] Buf;
        public int WPos;
        public int RPos;
        public int Capicity;

        public int Base;
        public bool isWriteGrow = false;

        public BytesStream()
        {
            Reset(null);
        }

        public BytesStream(byte[] bytes)
        {
            Reset(bytes);
        }

        public void SetWriteGrow(bool b)
        {
            isWriteGrow = b;
        }

        public void Reset()
        {
            this.Base = 0;
            this.RPos = 0;
            this.WPos = 0;
        }

        public void Reset(byte[] bytes, int inLen = -1)
        {
            this.Buf = bytes;
            this.Base = 0;
            this.RPos = 0;
            this.WPos = 0;
            if (bytes != null)
            {
                this.Capicity = (inLen == -1) ? bytes.Length : inLen;
            }
            else
            {
                this.Capicity = 0;
            }
        }

        public void Reset(byte[] bytes, int offset, int count)
        {
            this.Base = offset;
            this.Buf = bytes;
            this.RPos = offset;
            this.WPos = offset;
            this.Capicity = count + offset;
        }

        public void WriteBool(bool b)
        {
            if (!IsWEnough(1))
            {
                throw new EndOfStreamException();
            }

            this.Buf[WPos] = (byte)(b ? 1 : 0);
            this.WPos += 1;
        }

        public void WriteByte(byte b1)
        {
            if (!IsWEnough(1))
            {
                throw new EndOfStreamException();
            }

            this.Buf[WPos] = b1;
            this.WPos += 1;
        }

        public void WriteUInt32Array(UInt32[] p)
        {
            if (p == null)
            {
                WriteUInt16(0);
                return;
            }

            int len = p.Length;
            WriteUInt16((ushort)len);
            for (int i = 0; i < len; ++i)
            {
                WriteUInt32(p[i]);
            }
        }

        public void WriteBytes(byte[] p)
        {
            if (p == null)
            {
                WriteUInt32(0);
                return;
            }
            if (!IsWEnough(4+ p.Length))
            {
                throw new EndOfStreamException();
            }
            int len = p.Length;
            WriteUInt32((uint)len);
            WriteBuf(p, len);
        }

        public void WriteBuf(byte[] p, int len)
        {
            if (p == null || len == 0)
            {
                return;
            }

            if (!IsWEnough(len))
            {
                throw new EndOfStreamException();
            }

            Buffer.BlockCopy(p, 0, this.Buf, this.WPos, len);
            this.WPos += len;
        }

        public void WriteInt16(short n)
        {
            if (!IsWEnough(2))
            {
                throw new EndOfStreamException();
            }

            this.Buf[this.WPos + 0] = (byte)n;
            this.Buf[this.WPos + 1] = (byte)(n >> 8);
            this.WPos += 2;
        }

        public void WriteUInt16(ushort n)
        {
            if (!IsWEnough(2))
            {
                throw new EndOfStreamException();
            }

            this.Buf[this.WPos + 0] = (byte)n;
            this.Buf[this.WPos + 1] = (byte)(n >> 8);
            this.WPos += 2;
        }

        public void WriteInt32(int n)
        {
            if (!IsWEnough(4))
            {
                throw new EndOfStreamException();
            }

            this.Buf[this.WPos + 0] = (byte)(n & 0xff);
            this.Buf[this.WPos + 1] = (byte)((n >> 8) & 0xff);
            this.Buf[this.WPos + 2] = (byte)((n >> 16) & 0xff);
            this.Buf[this.WPos + 3] = (byte)((n >> 24) & 0xff);

            this.WPos += 4;
        }

        public void WriteUInt32(uint n)
        {
            if (!IsWEnough(4))
            {
                throw new EndOfStreamException();
            }

            this.Buf[this.WPos + 0] = (byte)(n & 0xff);
            this.Buf[this.WPos + 1] = (byte)((n >> 8) & 0xff);
            this.Buf[this.WPos + 2] = (byte)((n >> 16) & 0xff);
            this.Buf[this.WPos + 3] = (byte)((n >> 24) & 0xff);

            this.WPos += 4;
        }

        public void WriteInt64(long n)
        {
            WriteBuf(BitConverter.GetBytes(n), 8);
        }

        public void WriteUInt64(UInt64 n)
        {
            WriteBuf(BitConverter.GetBytes(n), 8);
        }

        public void WriteFloat(float n)
        {
            WriteBuf(BitConverter.GetBytes(n), 4);
        }

        public void WriteDouble(double n)
        {
            WriteBuf(BitConverter.GetBytes(n), 8);
        }

        public void WriteString(string str)
        {
            byte[] p = System.Text.Encoding.UTF8.GetBytes(str);
            if (p == null)
            {
                WriteUInt16(0);
                return;
            }

            int len = p.Length;
            WriteUInt16((UInt16)len);
            WriteBuf(p, len);
        }

        public byte ReadByte()
        {
            if (!IsREnough(1))
            {
                throw new EndOfStreamException();
            }

            this.RPos++;
            return this.Buf[this.RPos - 1];
        }

        public bool ReadBool()
        {
            if (!IsREnough(1))
            {
                throw new EndOfStreamException();
            }

            this.RPos++;
            return Convert.ToBoolean(this.Buf[this.RPos - 1]);
        }


        public byte[] ReadBuf(int count)
        {
            if (!IsREnough(count))
            {
                return null;
            }

            byte[] bytes = new byte[count];
            Buffer.BlockCopy(this.Buf, this.RPos, bytes, 0, count);
            this.RPos += count;
            return bytes;
        }

        public byte[] ReadBytes()
        {
            var len = ReadUInt32();
            return ReadBuf((int)len);
        }

        public uint[] ReadUInt32Array()
        {
            // Use ushort to persist the array length.
            var len = ReadUInt16();
            uint[] p = new uint[len];
            for (int i = 0; i < len; ++i)
            {
                p[i] = ReadUInt32();
            }

            return p;
        }

        public char ReadChar()
        {
            return (char)ReadByte();
        }

        public float ReadFloat()
        {
            if (!IsREnough(4))
            {
                throw new EndOfStreamException();
            }

            float ret = BitConverter.ToSingle(this.Buf, this.RPos);
            this.RPos += 4;
            return ret;
        }

        public double ReadDouble()
        {
            if (!IsREnough(8))
            {
                throw new EndOfStreamException();
            }

            double ret = BitConverter.ToDouble(this.Buf, this.RPos);
            this.RPos += 8;
            return ret;
        }

        public short ReadInt16()
        {
            if (!IsREnough(2))
            {
                throw new EndOfStreamException();
            }

            short ret = BitConverter.ToInt16(this.Buf, this.RPos);
            this.RPos += 2;
            return ret;
        }

        public ushort ReadUInt16()
        {
            if (!IsREnough(2))
            {
                throw new EndOfStreamException();
            }

            ushort ret = BitConverter.ToUInt16(this.Buf, this.RPos);
            this.RPos += 2;
            return ret;
        }

        public int ReadInt32()
        {
            if (!IsREnough(4))
            {
                throw new EndOfStreamException();
            }

            int ret = BitConverter.ToInt32(this.Buf, this.RPos);
            this.RPos += 4;
            return ret;
        }

        public uint ReadUInt32()
        {
            if (!IsREnough(4))
            {
                return 0;
            }

            uint ret = BitConverter.ToUInt32(this.Buf, this.RPos);
            this.RPos += 4;
            return ret;
        }

        public long ReadInt64()
        {
            if (!IsREnough(8))
            {
                throw new EndOfStreamException();
            }

            long ret = BitConverter.ToInt64(this.Buf, this.RPos);
            this.RPos += 8;
            return ret;
        }

        public ulong ReadUInt64()
        {
            if (!IsREnough(8))
            {
                return 0;
            }

            ulong ret = BitConverter.ToUInt64(this.Buf, this.RPos);
            this.RPos += 8;
            return ret;
        }

        public string ReadString()
        {
            int len = (int)ReadUInt16();

            if (!IsREnough(len))
            {
                throw new EndOfStreamException();
            }

            string ret = System.Text.Encoding.UTF8.GetString(this.Buf, this.RPos, len);
            this.RPos += len;

            return ret;
        }

        public bool IsWEnough(int size)
        {
            if (this.isWriteGrow)
            {
                if (this.Capicity < (size + this.WPos))
                {
                    this.growByteBuffer(size + this.WPos);
                }
            }
      
            return (this.Capicity - this.WPos) >= size;
        }

        public bool IsREnough(int size)
        {
            return (this.Capicity - this.RPos) >= size;
        }
        public bool IsReadEnd()
        {
            return this.Capicity == this.RPos;
        }

        public bool IsEOF()
        {
            return this.RPos == this.Capicity;
        }

        public bool growByteBuffer(int needSize)
        {
            var newLen = 0;
            if (this.Capicity>=MAX_STREAM_SIZE || needSize >= MAX_STREAM_SIZE) {
                return false;
            } 

            if (this.Capicity> MAX_STREAM_SIZE/2 || needSize > MAX_STREAM_SIZE / 2)
            {
                newLen = MAX_STREAM_SIZE - this.Capicity;
            } 
            else
            {
                newLen = Math.Max(STREAM_GROW_MIN_SIZE, this.Capicity);
                if (needSize > newLen + this.Capicity)
                {
                    newLen = needSize;
                }
            }

            var extension = new byte[newLen+this.Capicity];
            Buffer.BlockCopy(this.Buf, 0, extension, 0, this.Capicity);
            this.Capicity = extension.Length;
            this.Buf = extension;
            return true;
        }


        public const int MAX_STREAM_SIZE = 8 * 1024 * 1024;//8MB
        public const int STREAM_GROW_MIN_SIZE = 128;    //一次增长最小

        public static int GetStreamBaseTypeSize(object o,out bool found)
        {
            found = false;
            var t = o.GetType();
            
            if (t.Equals(typeof(byte))|| t.Equals(typeof(bool)) || t.Equals(typeof(char)))
            {
                found = true;
                return 1;
            }

            if (t.Equals(typeof(Int16)) || t.Equals(typeof(UInt16)))
            {
                found = true;
                return 2;
            }

            if (t.Equals(typeof(int)) || t.Equals(typeof(Int32)) || t.Equals(typeof(UInt32)) || t.Equals(typeof(float)))
            {
                found = true;
                return 4;
            }

            if (t.Equals(typeof(Int64)) || t.Equals(typeof(UInt64)) || t.Equals(typeof(double)))
            {
                found = true;
                return 8;
            }

            if (t.Equals(typeof(string)))
            {
                found = true;
                byte[] p = Encoding.UTF8.GetBytes(o as string);
                return 2+p.Length;
            }

            if (t.Equals(typeof(byte[])))
            {
                found = true;
                byte[] p = (byte[])o;
                return 4 + p.Length;
            }

            return 0;
        }
    }

}

