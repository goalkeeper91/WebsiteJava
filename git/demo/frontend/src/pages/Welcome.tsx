import { Link } from "react-router-dom";
import { CheckCircle2 } from "lucide-react";
import Seo from "../components/Seo";

// Paddle's overlay checkout redirects here via settings.successUrl once a
// transaction completes - the actual tier upgrade only lands a little later
// via the Paddle webhook, so this page deliberately doesn't try to read/show
// the new subscription state itself (avoids a "still shows old tier" flash).
export default function Welcome() {
  return (
    <div className="relative w-full min-h-screen bg-slate-950 text-white flex items-center justify-center px-6">
      <Seo
        title="Willkommen an Bord"
        description="Danke für dein Abo bei Goalkeeper91."
        path="/welcome"
        noindex
      />
      <div className="text-center max-w-md">
        <CheckCircle2 className="w-16 h-16 text-green-500 mx-auto mb-6" />
        <h1 className="text-3xl font-bold mb-4">Willkommen an Bord!</h1>
        <p className="text-gray-300 mb-8">
          Danke für dein Abo. Es kann einen kurzen Moment dauern, bis dein neuer Tarif im Dashboard
          angezeigt wird.
        </p>
        <Link
          to="/dashboard"
          className="inline-block py-3 px-6 bg-green-600 hover:bg-green-500 rounded-lg transition-colors font-semibold"
        >
          Zum Dashboard
        </Link>
      </div>
    </div>
  );
}
