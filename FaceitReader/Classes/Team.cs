using System;

namespace FaceitReader.Classes
{
    public class Team : IEquatable<Team>
    {
        public int Id { get; set; }
        public int Place { get; set; }
        public string TeamName { get; set; }
        public string TeamId { get; set; }
        public string CaptainName { get; set; }
        public override string ToString()
        {
            return "Place: " + Place + "    TeamID: " + TeamId + "  TeamName: " + TeamName + "   CaptainName: " + CaptainName;
        }
        public override bool Equals(object obj)
        {
            if (obj == null) return false;
            Team objAsTeam = obj as Team;
            if (objAsTeam == null) return false;
            else return Equals(objAsTeam);
        }
        public override int GetHashCode()
        {
            return Id;
        }
        public bool Equals(Team other)
        {
            if (other == null) return false;
            return (this.Id.Equals(other.Id));
        }
    }
}
