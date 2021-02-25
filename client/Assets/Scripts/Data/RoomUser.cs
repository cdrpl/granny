namespace Idlemon.Data
{
    /// <summary>
    /// Represents a user in a lobby room.
    /// </summary>
    public class RoomUser
    {
        public int Id { get; set; }
        public string Name { get; set; }
        public bool IsReady { get; set; }
    }
}
