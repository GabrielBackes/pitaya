// <auto-generated>
//     Generated by the protocol buffer compiler.  DO NOT EDIT!
//     source: msg.proto
// </auto-generated>
#pragma warning disable 1591, 0612, 3021
#region Designer generated code

using pb = global::Google.Protobuf;
using pbc = global::Google.Protobuf.Collections;
using pbr = global::Google.Protobuf.Reflection;
using scg = global::System.Collections.Generic;
namespace NPitaya.Protos {

  /// <summary>Holder for reflection information generated from msg.proto</summary>
  public static partial class MsgReflection {

    #region Descriptor
    /// <summary>File descriptor for msg.proto</summary>
    public static pbr::FileDescriptor Descriptor {
      get { return descriptor; }
    }
    private static pbr::FileDescriptor descriptor;

    static MsgReflection() {
      byte[] descriptorData = global::System.Convert.FromBase64String(
          string.Concat(
            "Cgltc2cucHJvdG8SBnByb3RvcyJcCgNNc2cSCgoCaWQYASABKAQSDQoFcm91",
            "dGUYAiABKAkSDAoEZGF0YRgDIAEoDBINCgVyZXBseRgEIAEoCRIdCgR0eXBl",
            "GAUgASgOMg8ucHJvdG9zLk1zZ1R5cGUqRgoHTXNnVHlwZRIOCgpNc2dSZXF1",
            "ZXN0EAASDQoJTXNnTm90aWZ5EAESDwoLTXNnUmVzcG9uc2UQAhILCgdNc2dQ",
            "dXNoEANCPFopZ2l0aHViLmNvbS90b3BmcmVlZ2FtZXMvcGl0YXlhL3BrZy9w",
            "cm90b3OqAg5OUGl0YXlhLlByb3Rvc2IGcHJvdG8z"));
      descriptor = pbr::FileDescriptor.FromGeneratedCode(descriptorData,
          new pbr::FileDescriptor[] { },
          new pbr::GeneratedClrTypeInfo(new[] {typeof(global::NPitaya.Protos.MsgType), }, null, new pbr::GeneratedClrTypeInfo[] {
            new pbr::GeneratedClrTypeInfo(typeof(global::NPitaya.Protos.Msg), global::NPitaya.Protos.Msg.Parser, new[]{ "Id", "Route", "Data", "Reply", "Type" }, null, null, null, null)
          }));
    }
    #endregion

  }
  #region Enums
  public enum MsgType {
    [pbr::OriginalName("MsgRequest")] MsgRequest = 0,
    [pbr::OriginalName("MsgNotify")] MsgNotify = 1,
    [pbr::OriginalName("MsgResponse")] MsgResponse = 2,
    [pbr::OriginalName("MsgPush")] MsgPush = 3,
  }

  #endregion

  #region Messages
  public sealed partial class Msg : pb::IMessage<Msg>
  #if !GOOGLE_PROTOBUF_REFSTRUCT_COMPATIBILITY_MODE
      , pb::IBufferMessage
  #endif
  {
    private static readonly pb::MessageParser<Msg> _parser = new pb::MessageParser<Msg>(() => new Msg());
    private pb::UnknownFieldSet _unknownFields;
    [global::System.Diagnostics.DebuggerNonUserCodeAttribute]
    public static pb::MessageParser<Msg> Parser { get { return _parser; } }

    [global::System.Diagnostics.DebuggerNonUserCodeAttribute]
    public static pbr::MessageDescriptor Descriptor {
      get { return global::NPitaya.Protos.MsgReflection.Descriptor.MessageTypes[0]; }
    }

    [global::System.Diagnostics.DebuggerNonUserCodeAttribute]
    pbr::MessageDescriptor pb::IMessage.Descriptor {
      get { return Descriptor; }
    }

    [global::System.Diagnostics.DebuggerNonUserCodeAttribute]
    public Msg() {
      OnConstruction();
    }

    partial void OnConstruction();

    [global::System.Diagnostics.DebuggerNonUserCodeAttribute]
    public Msg(Msg other) : this() {
      id_ = other.id_;
      route_ = other.route_;
      data_ = other.data_;
      reply_ = other.reply_;
      type_ = other.type_;
      _unknownFields = pb::UnknownFieldSet.Clone(other._unknownFields);
    }

    [global::System.Diagnostics.DebuggerNonUserCodeAttribute]
    public Msg Clone() {
      return new Msg(this);
    }

    /// <summary>Field number for the "id" field.</summary>
    public const int IdFieldNumber = 1;
    private ulong id_;
    [global::System.Diagnostics.DebuggerNonUserCodeAttribute]
    public ulong Id {
      get { return id_; }
      set {
        id_ = value;
      }
    }

    /// <summary>Field number for the "route" field.</summary>
    public const int RouteFieldNumber = 2;
    private string route_ = "";
    [global::System.Diagnostics.DebuggerNonUserCodeAttribute]
    public string Route {
      get { return route_; }
      set {
        route_ = pb::ProtoPreconditions.CheckNotNull(value, "value");
      }
    }

    /// <summary>Field number for the "data" field.</summary>
    public const int DataFieldNumber = 3;
    private pb::ByteString data_ = pb::ByteString.Empty;
    [global::System.Diagnostics.DebuggerNonUserCodeAttribute]
    public pb::ByteString Data {
      get { return data_; }
      set {
        data_ = pb::ProtoPreconditions.CheckNotNull(value, "value");
      }
    }

    /// <summary>Field number for the "reply" field.</summary>
    public const int ReplyFieldNumber = 4;
    private string reply_ = "";
    [global::System.Diagnostics.DebuggerNonUserCodeAttribute]
    public string Reply {
      get { return reply_; }
      set {
        reply_ = pb::ProtoPreconditions.CheckNotNull(value, "value");
      }
    }

    /// <summary>Field number for the "type" field.</summary>
    public const int TypeFieldNumber = 5;
    private global::NPitaya.Protos.MsgType type_ = global::NPitaya.Protos.MsgType.MsgRequest;
    [global::System.Diagnostics.DebuggerNonUserCodeAttribute]
    public global::NPitaya.Protos.MsgType Type {
      get { return type_; }
      set {
        type_ = value;
      }
    }

    [global::System.Diagnostics.DebuggerNonUserCodeAttribute]
    public override bool Equals(object other) {
      return Equals(other as Msg);
    }

    [global::System.Diagnostics.DebuggerNonUserCodeAttribute]
    public bool Equals(Msg other) {
      if (ReferenceEquals(other, null)) {
        return false;
      }
      if (ReferenceEquals(other, this)) {
        return true;
      }
      if (Id != other.Id) return false;
      if (Route != other.Route) return false;
      if (Data != other.Data) return false;
      if (Reply != other.Reply) return false;
      if (Type != other.Type) return false;
      return Equals(_unknownFields, other._unknownFields);
    }

    [global::System.Diagnostics.DebuggerNonUserCodeAttribute]
    public override int GetHashCode() {
      int hash = 1;
      if (Id != 0UL) hash ^= Id.GetHashCode();
      if (Route.Length != 0) hash ^= Route.GetHashCode();
      if (Data.Length != 0) hash ^= Data.GetHashCode();
      if (Reply.Length != 0) hash ^= Reply.GetHashCode();
      if (Type != global::NPitaya.Protos.MsgType.MsgRequest) hash ^= Type.GetHashCode();
      if (_unknownFields != null) {
        hash ^= _unknownFields.GetHashCode();
      }
      return hash;
    }

    [global::System.Diagnostics.DebuggerNonUserCodeAttribute]
    public override string ToString() {
      return pb::JsonFormatter.ToDiagnosticString(this);
    }

    [global::System.Diagnostics.DebuggerNonUserCodeAttribute]
    public void WriteTo(pb::CodedOutputStream output) {
    #if !GOOGLE_PROTOBUF_REFSTRUCT_COMPATIBILITY_MODE
      output.WriteRawMessage(this);
    #else
      if (Id != 0UL) {
        output.WriteRawTag(8);
        output.WriteUInt64(Id);
      }
      if (Route.Length != 0) {
        output.WriteRawTag(18);
        output.WriteString(Route);
      }
      if (Data.Length != 0) {
        output.WriteRawTag(26);
        output.WriteBytes(Data);
      }
      if (Reply.Length != 0) {
        output.WriteRawTag(34);
        output.WriteString(Reply);
      }
      if (Type != global::NPitaya.Protos.MsgType.MsgRequest) {
        output.WriteRawTag(40);
        output.WriteEnum((int) Type);
      }
      if (_unknownFields != null) {
        _unknownFields.WriteTo(output);
      }
    #endif
    }

    #if !GOOGLE_PROTOBUF_REFSTRUCT_COMPATIBILITY_MODE
    [global::System.Diagnostics.DebuggerNonUserCodeAttribute]
    void pb::IBufferMessage.InternalWriteTo(ref pb::WriteContext output) {
      if (Id != 0UL) {
        output.WriteRawTag(8);
        output.WriteUInt64(Id);
      }
      if (Route.Length != 0) {
        output.WriteRawTag(18);
        output.WriteString(Route);
      }
      if (Data.Length != 0) {
        output.WriteRawTag(26);
        output.WriteBytes(Data);
      }
      if (Reply.Length != 0) {
        output.WriteRawTag(34);
        output.WriteString(Reply);
      }
      if (Type != global::NPitaya.Protos.MsgType.MsgRequest) {
        output.WriteRawTag(40);
        output.WriteEnum((int) Type);
      }
      if (_unknownFields != null) {
        _unknownFields.WriteTo(ref output);
      }
    }
    #endif

    [global::System.Diagnostics.DebuggerNonUserCodeAttribute]
    public int CalculateSize() {
      int size = 0;
      if (Id != 0UL) {
        size += 1 + pb::CodedOutputStream.ComputeUInt64Size(Id);
      }
      if (Route.Length != 0) {
        size += 1 + pb::CodedOutputStream.ComputeStringSize(Route);
      }
      if (Data.Length != 0) {
        size += 1 + pb::CodedOutputStream.ComputeBytesSize(Data);
      }
      if (Reply.Length != 0) {
        size += 1 + pb::CodedOutputStream.ComputeStringSize(Reply);
      }
      if (Type != global::NPitaya.Protos.MsgType.MsgRequest) {
        size += 1 + pb::CodedOutputStream.ComputeEnumSize((int) Type);
      }
      if (_unknownFields != null) {
        size += _unknownFields.CalculateSize();
      }
      return size;
    }

    [global::System.Diagnostics.DebuggerNonUserCodeAttribute]
    public void MergeFrom(Msg other) {
      if (other == null) {
        return;
      }
      if (other.Id != 0UL) {
        Id = other.Id;
      }
      if (other.Route.Length != 0) {
        Route = other.Route;
      }
      if (other.Data.Length != 0) {
        Data = other.Data;
      }
      if (other.Reply.Length != 0) {
        Reply = other.Reply;
      }
      if (other.Type != global::NPitaya.Protos.MsgType.MsgRequest) {
        Type = other.Type;
      }
      _unknownFields = pb::UnknownFieldSet.MergeFrom(_unknownFields, other._unknownFields);
    }

    [global::System.Diagnostics.DebuggerNonUserCodeAttribute]
    public void MergeFrom(pb::CodedInputStream input) {
    #if !GOOGLE_PROTOBUF_REFSTRUCT_COMPATIBILITY_MODE
      input.ReadRawMessage(this);
    #else
      uint tag;
      while ((tag = input.ReadTag()) != 0) {
        switch(tag) {
          default:
            _unknownFields = pb::UnknownFieldSet.MergeFieldFrom(_unknownFields, input);
            break;
          case 8: {
            Id = input.ReadUInt64();
            break;
          }
          case 18: {
            Route = input.ReadString();
            break;
          }
          case 26: {
            Data = input.ReadBytes();
            break;
          }
          case 34: {
            Reply = input.ReadString();
            break;
          }
          case 40: {
            Type = (global::NPitaya.Protos.MsgType) input.ReadEnum();
            break;
          }
        }
      }
    #endif
    }

    #if !GOOGLE_PROTOBUF_REFSTRUCT_COMPATIBILITY_MODE
    [global::System.Diagnostics.DebuggerNonUserCodeAttribute]
    void pb::IBufferMessage.InternalMergeFrom(ref pb::ParseContext input) {
      uint tag;
      while ((tag = input.ReadTag()) != 0) {
        switch(tag) {
          default:
            _unknownFields = pb::UnknownFieldSet.MergeFieldFrom(_unknownFields, ref input);
            break;
          case 8: {
            Id = input.ReadUInt64();
            break;
          }
          case 18: {
            Route = input.ReadString();
            break;
          }
          case 26: {
            Data = input.ReadBytes();
            break;
          }
          case 34: {
            Reply = input.ReadString();
            break;
          }
          case 40: {
            Type = (global::NPitaya.Protos.MsgType) input.ReadEnum();
            break;
          }
        }
      }
    }
    #endif

  }

  #endregion

}

#endregion Designer generated code
