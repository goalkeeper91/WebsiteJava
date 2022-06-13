using System;

namespace FaceitReader.Classes
{
    public class Member : IEquatable<Member>
    {
        public int Id { get; set; }
        public string PlayerId { get; set; }
        public string PlayerName { get; set; }
        public string TeamId { get; set; }
        public override string ToString()
        {
            return "TeamId: " + TeamId + "PlayerId: " + PlayerId + "  PlayerName: " + PlayerName;
        }
        public override bool Equals(object obj)
        {
            if (obj == null) return false;
            Member objAsMember = obj as Member;
            if (objAsMember == null) return false;
            else return Equals(objAsMember);
        }
        public override int GetHashCode()
        {
            return Id;
        }
        public bool Equals(Member other)
        {
            if (other == null) return false;
            return (this.Id.Equals(other.Id));
        }
    }
}
