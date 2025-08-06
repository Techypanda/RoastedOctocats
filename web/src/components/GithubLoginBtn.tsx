import * as React from "react";
import { Button, makeStyles, shorthands, tokens } from "@fluentui/react-components";
import { v7 } from "uuid";

// SVG for GitHub logo (24px, regular style)
const GithubIcon = () => (
  <svg
    width="24"
    height="24"
    viewBox="0 0 24 24"
    fill="none"
    aria-hidden="true"
    focusable="false"
    xmlns="http://www.w3.org/2000/svg"
  >
    <path
      fill="#fff"
      d="M12 2C6.477 2 2 6.484 2 12.012c0 4.429 2.868 8.184 6.839 9.521.5.093.682-.217.682-.483 0-.237-.009-.868-.013-1.703-2.782.605-3.369-1.342-3.369-1.342-.454-1.155-1.11-1.463-1.11-1.463-.908-.62.069-.608.069-.608 1.004.07 1.532 1.032 1.532 1.032.892 1.529 2.341 1.088 2.91.833.092-.647.349-1.088.636-1.34-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.03-2.688-.103-.253-.447-1.27.098-2.647 0 0 .84-.27 2.75 1.025A9.564 9.564 0 0 1 12 6.844c.85.004 1.705.115 2.504.337 1.909-1.295 2.747-1.025 2.747-1.025.546 1.377.202 2.394.1 2.647.64.7 1.03 1.595 1.03 2.688 0 3.848-2.338 4.695-4.566 4.944.359.309.68.92.68 1.854 0 1.338-.012 2.419-.012 2.748 0 .268.18.58.688.482C19.135 20.192 22 16.44 22 12.012 22 6.484 17.523 2 12 2Z"
    />
  </svg>
);

const useStyles = makeStyles({
  root: {
    ...shorthands.padding("8px", "16px"),
    ...shorthands.borderRadius(tokens.borderRadiusMedium),
    backgroundColor: "#24292f",
    color: "#fff",
    border: "none",
    display: "inline-flex",
    alignItems: "center",
    gap: "8px",
    fontWeight: 500,
    boxShadow: tokens.shadow4,
    cursor: "pointer",
    transition: "background 0.2s",
    ":hover": {
      backgroundColor: "#1b1f23",
    },
    ":active": {
      backgroundColor: "#16181c",
    },
  },
  githubIcon: {
    display: "inline-block",
    verticalAlign: "middle",
  }
});

async function generateCodeChallenge(codeVerifier: string) {
    var digest = await crypto.subtle.digest("SHA-256",
        new TextEncoder().encode(codeVerifier));

    return btoa(String.fromCharCode(...new Uint8Array(digest)))
        .replace(/=/g, '').replace(/\+/g, '-').replace(/\//g, '_')
}

function generateRandomString(length: number = 64) {
    var text = "";
    var possible = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";

    for (var i = 0; i < length; i++) {
        text += possible.charAt(Math.floor(Math.random() * possible.length));
    }

    return text;
}

const loginWithGithub = async () => {
    const pkceStringGenerated = generateRandomString();
    const shad = await generateCodeChallenge(pkceStringGenerated);
    const state = v7();
    const url = new URL('https://github.com/login/oauth/authorize');
    url.searchParams.append('client_id', 'Iv23ligun1uyOZYdvxnq');
    url.searchParams.append('redirect_uri', `${window.location.origin}/postlogin`);
    url.searchParams.append('scope', '');
    url.searchParams.append('state', state)
    url.searchParams.append('code_challenge', shad)
    url.searchParams.append('code_challenge_method', 'S256')
    url.searchParams.append('allow_signup', 'true')
    // save it to session storage
    sessionStorage.setItem(state, JSON.stringify({ original: pkceStringGenerated, sha: shad }));
    window.location.href = url.toString();
}

export const GithubLoginButton: React.FC<React.ButtonHTMLAttributes<HTMLButtonElement>> = (props) => {
  const styles = useStyles();
  return (
    <Button
      className={styles.root}
      icon={<span className={styles.githubIcon}><GithubIcon /></span>}
      {...props}
      onClick={() => loginWithGithub()}
    >
      Login with GitHub
    </Button>
  );
};