import React from "react";
import { renderToString } from "react-dom/server";
import { createGlobalStyle, css, ServerStyleSheet } from "styled-components";

React.useLayoutEffect = React.useEffect;

const GlobalStyle = createGlobalStyle`${css`
  div {
    color: red;
  }
`}`;

export default () => {
  const sheet = new ServerStyleSheet();

  const html = renderToString(
    sheet.collectStyles(
      <>
        <GlobalStyle />
        <div>Hello SSR</div>
      </>
    )
  );

  const styleTags = sheet.getStyleTags();

  sheet.seal();

  return `${styleTags}${html}`;
};
