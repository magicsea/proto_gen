
using System;

namespace Misc
{
    public interface IMessage
    {
        UInt16 MsgID { get; }
        string MsgName { get; }
        void Write(BytesStream s);
        void Read(BytesStream s);
        int GetStreamSize();
    }

    public interface IPBMessage: Google.Protobuf.IMessage
    {
        UInt16 MsgID { get; }
        string MsgName { get; }

    }
}