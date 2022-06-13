using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace FaceitReader.Classes
{
    public class Placement : IEquatable<Placement>
    {
        public int Id { get; set; }
        public int Place { get; set; }
        public string TeamId { get; set; }
        public string TeamName { get; set; }
        public override string ToString()
        {
            return "Place: " + Place + "TeamId: " + TeamId + "  Name: " + TeamName;
        }
        public override bool Equals(object obj)
        {
            if (obj == null) return false;
            Placement objAsPlacement = obj as Placement;
            if (objAsPlacement == null) return false;
            else return Equals(objAsPlacement);
        }
        public override int GetHashCode()
        {
            return Id;
        }
        public bool Equals(Placement other)
        {
            if (other == null) return false;
            return (this.Id.Equals(other.Id));
        }
    }
}
