// <auto-generated>
//     Generated by the protocol buffer compiler.  DO NOT EDIT!
//     source: granny.proto
// </auto-generated>
#pragma warning disable 0414, 1591
#region Designer generated code

using grpc = global::Grpc.Core;

namespace Proto {
  public static partial class Auth
  {
    static readonly string __ServiceName = "proto.Auth";

    static readonly grpc::Marshaller<global::Proto.SignUpRequest> __Marshaller_proto_SignUpRequest = grpc::Marshallers.Create((arg) => global::Google.Protobuf.MessageExtensions.ToByteArray(arg), global::Proto.SignUpRequest.Parser.ParseFrom);
    static readonly grpc::Marshaller<global::Proto.SignUpResponse> __Marshaller_proto_SignUpResponse = grpc::Marshallers.Create((arg) => global::Google.Protobuf.MessageExtensions.ToByteArray(arg), global::Proto.SignUpResponse.Parser.ParseFrom);

    static readonly grpc::Method<global::Proto.SignUpRequest, global::Proto.SignUpResponse> __Method_SignUp = new grpc::Method<global::Proto.SignUpRequest, global::Proto.SignUpResponse>(
        grpc::MethodType.Unary,
        __ServiceName,
        "SignUp",
        __Marshaller_proto_SignUpRequest,
        __Marshaller_proto_SignUpResponse);

    /// <summary>Service descriptor</summary>
    public static global::Google.Protobuf.Reflection.ServiceDescriptor Descriptor
    {
      get { return global::Proto.GrannyReflection.Descriptor.Services[0]; }
    }

    /// <summary>Base class for server-side implementations of Auth</summary>
    [grpc::BindServiceMethod(typeof(Auth), "BindService")]
    public abstract partial class AuthBase
    {
      public virtual global::System.Threading.Tasks.Task<global::Proto.SignUpResponse> SignUp(global::Proto.SignUpRequest request, grpc::ServerCallContext context)
      {
        throw new grpc::RpcException(new grpc::Status(grpc::StatusCode.Unimplemented, ""));
      }

    }

    /// <summary>Client for Auth</summary>
    public partial class AuthClient : grpc::ClientBase<AuthClient>
    {
      /// <summary>Creates a new client for Auth</summary>
      /// <param name="channel">The channel to use to make remote calls.</param>
      public AuthClient(grpc::ChannelBase channel) : base(channel)
      {
      }
      /// <summary>Creates a new client for Auth that uses a custom <c>CallInvoker</c>.</summary>
      /// <param name="callInvoker">The callInvoker to use to make remote calls.</param>
      public AuthClient(grpc::CallInvoker callInvoker) : base(callInvoker)
      {
      }
      /// <summary>Protected parameterless constructor to allow creation of test doubles.</summary>
      protected AuthClient() : base()
      {
      }
      /// <summary>Protected constructor to allow creation of configured clients.</summary>
      /// <param name="configuration">The client configuration.</param>
      protected AuthClient(ClientBaseConfiguration configuration) : base(configuration)
      {
      }

      public virtual global::Proto.SignUpResponse SignUp(global::Proto.SignUpRequest request, grpc::Metadata headers = null, global::System.DateTime? deadline = null, global::System.Threading.CancellationToken cancellationToken = default(global::System.Threading.CancellationToken))
      {
        return SignUp(request, new grpc::CallOptions(headers, deadline, cancellationToken));
      }
      public virtual global::Proto.SignUpResponse SignUp(global::Proto.SignUpRequest request, grpc::CallOptions options)
      {
        return CallInvoker.BlockingUnaryCall(__Method_SignUp, null, options, request);
      }
      public virtual grpc::AsyncUnaryCall<global::Proto.SignUpResponse> SignUpAsync(global::Proto.SignUpRequest request, grpc::Metadata headers = null, global::System.DateTime? deadline = null, global::System.Threading.CancellationToken cancellationToken = default(global::System.Threading.CancellationToken))
      {
        return SignUpAsync(request, new grpc::CallOptions(headers, deadline, cancellationToken));
      }
      public virtual grpc::AsyncUnaryCall<global::Proto.SignUpResponse> SignUpAsync(global::Proto.SignUpRequest request, grpc::CallOptions options)
      {
        return CallInvoker.AsyncUnaryCall(__Method_SignUp, null, options, request);
      }
      /// <summary>Creates a new instance of client from given <c>ClientBaseConfiguration</c>.</summary>
      protected override AuthClient NewInstance(ClientBaseConfiguration configuration)
      {
        return new AuthClient(configuration);
      }
    }

    /// <summary>Creates service definition that can be registered with a server</summary>
    /// <param name="serviceImpl">An object implementing the server-side handling logic.</param>
    public static grpc::ServerServiceDefinition BindService(AuthBase serviceImpl)
    {
      return grpc::ServerServiceDefinition.CreateBuilder()
          .AddMethod(__Method_SignUp, serviceImpl.SignUp).Build();
    }

    /// <summary>Register service method with a service binder with or without implementation. Useful when customizing the  service binding logic.
    /// Note: this method is part of an experimental API that can change or be removed without any prior notice.</summary>
    /// <param name="serviceBinder">Service methods will be bound by calling <c>AddMethod</c> on this object.</param>
    /// <param name="serviceImpl">An object implementing the server-side handling logic.</param>
    public static void BindService(grpc::ServiceBinderBase serviceBinder, AuthBase serviceImpl)
    {
      serviceBinder.AddMethod(__Method_SignUp, serviceImpl == null ? null : new grpc::UnaryServerMethod<global::Proto.SignUpRequest, global::Proto.SignUpResponse>(serviceImpl.SignUp));
    }

  }
}
#endregion