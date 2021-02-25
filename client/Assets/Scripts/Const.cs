namespace Idlemon
{
    /// <summary>
    /// Constant values.
    /// </summary>
    public static class Const
    {
        /// <summary>
        /// The web protocol used to access the APIs.
        /// </summary>
        public const string WEB_PROTOCOL = "http";

        /// <summary>
        /// The address of the public facing API.
        /// </summary>
        public const string PUBLIC_API = "127.0.0.1:3000";

        /// <summary>
        /// The address of the lobby server.
        /// </summary>
        public const string LOBBY_ADDR = "127.0.0.1:3001";

        /// <summary>
        /// TCP messages are sent over channels.
        /// </summary>
        public enum Channel
        {
            /// <summary> 
            /// For authenticating clients.
            /// </summary>
            AuthRequest,

            /// <summary> 
            /// Response to clients after receiving an Auth request.
            /// </summary>
            AuthResponse,

            /// <summary> 
            /// Clients are pinged periodically to detect disconnections.
            /// </summary>
            Ping,

            /// <summary> 
            /// Chat messages are sent through this channel.
            /// </summary>
            Chat,
        }
    }
}
