using System;

namespace FaceitReader.Classes
{
    public class Captain : IEquatable<Captain>
    {
        public int Id { get; set; }
        public string TeamId { get; set; }
        public string CaptainId { get; set; }
        public override string ToString()
        {
            return "TeamId: " + TeamId + "CaptainId: " + CaptainId;
        }
        public override bool Equals(object obj)
        {
            if (obj == null) return false;
            Captain objAsCaptain = obj as Captain;
            if (objAsCaptain == null) return false;
            else return Equals(objAsCaptain);
        }
        public override int GetHashCode()
        {
            return Id;
        }
        public bool Equals(Captain other)
        {
            if (other == null) return false;
            return (this.Id.Equals(other.Id));
        }
    }
}
