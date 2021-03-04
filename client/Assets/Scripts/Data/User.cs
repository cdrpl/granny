namespace Idlemon.Data
{
    public class User
    {
        public int Id { get; set; }
        public string Name { get; set; }
        public string Token { get; set; }

        public User(Proto.SignInResponse response)
        {
            Id = response.Id;
            Name = response.Name;
            Token = response.Token;
        }
    }
}
