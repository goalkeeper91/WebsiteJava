package streamer_website.demo.handler;

import jakarta.servlet.*;
import jakarta.servlet.http.*;
import org.springframework.stereotype.Component;
import java.io.IOException;
import java.util.Arrays;
import java.util.List;

@Component
public class CookieConsentFilter implements Filter {

    // Kategorien definieren
    private final List<String> necessaryCookies = Arrays.asList(
            "JSESSIONID", "adminer_key", "adminer_permanent", "adminer_sid", "adminer_version"
    );

    private final List<String> statisticsCookies = Arrays.asList(
            "_ga","_gid","_gac_UA-","_ga_0WK9KVV1G4","_ga_4C7YTQZ45L","_ga_HXY5PN80YH",
            "_ga_WX15RJ0HZZ","_gcl_au","_gcl_aw","_gcl_gs","_hjSessionUser_","_BEAMER_USER_ID_"
    );

    private final List<String> marketingCookies = Arrays.asList(
            "_fbp","_uetsid","_uetvid","cb_mktg","AMP_MKTG_7041823e70","FPID","FPLC",
            "intercom-device-id-","intercom-id-","IR_25071","IR_gbd","IR_PI","ak_bmsc"
    );

    @Override
    public void doFilter(ServletRequest request, ServletResponse response, FilterChain chain)
            throws IOException, ServletException {

        HttpServletRequest req = (HttpServletRequest) request;
        HttpServletResponse res = (HttpServletResponse) response;

        // Consent prüfen
        boolean consentStatistics = false;
        boolean consentMarketing = false;

        Cookie[] cookies = req.getCookies();
        if (cookies != null) {
            for (Cookie c : cookies) {
                if ("CookieConsent".equals(c.getName())) {
                    String val = c.getValue();
                    consentStatistics = val.contains("statistics:true");
                    consentMarketing = val.contains("marketing:true");
                }
            }
        }

        // Notwendige Cookies immer setzen
        necessaryCookies.forEach(name -> {
            Cookie cookie = new Cookie(name, "example_value"); // ggf. dynamisch
            cookie.setPath("/");
            cookie.setHttpOnly(true);
            cookie.setMaxAge(3600); // Beispiel: 1 Stunde
            res.addCookie(cookie);
        });

        // Statistik-Cookies nur bei Consent
        if (consentStatistics) {
            statisticsCookies.forEach(name -> {
                Cookie cookie = new Cookie(name, "example_value");
                cookie.setPath("/");
                cookie.setHttpOnly(true);
                cookie.setMaxAge(365*24*60*60); // 1 Jahr
                res.addCookie(cookie);
            });
        }

        // Marketing-Cookies nur bei Consent
        if (consentMarketing) {
            marketingCookies.forEach(name -> {
                Cookie cookie = new Cookie(name, "example_value");
                cookie.setPath("/");
                cookie.setHttpOnly(true);
                cookie.setMaxAge(180*24*60*60); // 6 Monate
                res.addCookie(cookie);
            });
        }

        chain.doFilter(req, res);
    }
}
