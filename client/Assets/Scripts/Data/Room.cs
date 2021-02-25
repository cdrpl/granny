using System.Collections.Generic;

namespace Idlemon.Data
{
    /// <summary>
    /// Represents a lobby room.
    /// </summary>
    public class Room
    {
        public string Id { get; set; }
        public string Name { get; set; }
        public Dictionary<int, RoomUser> users { get; set; }

        /// <summary>
        /// The number of users currently in the room.
        /// </summary>
        public int NumUsers => users.Count;
    }
}
