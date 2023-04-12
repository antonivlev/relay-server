const round = (num, precision = 2) => {
  const factor = Math.pow(10, precision);
  return Math.round(num * factor) / factor;
};

const parseJson = jsonStr => {
  try {
    return JSON.parse(jsonStr);
  } catch (e) {
    return null;
  }
};

const getEmailFromAccessTokenFromCookies = () => {
  if (!document.cookie) {
    return null;
  }

  const accessToken = document.cookie
    .split(";")
    .find(row => row.startsWith("accessToken"))
    .split("accessToken=")?.[1];

  if (!accessToken) {
    return null;
  }

  const email = atob(accessToken)?.split(":")?.[0];
  return email;
};

const getAuthHeaderFromCookies = () => {
  if (!document.cookie) {
    return {};
  }

  const accessToken = document.cookie
    .split(";")
    .find(row => row.startsWith("accessToken"))
    .split("accessToken=")?.[1];

  if (!accessToken) {
    return {};
  }

  return {
    Authorization: `Basic ${accessToken}`,
  };
};

const doGet = async url => {
  const res = await fetch(url, {
    headers: {
      "Content-Type": "application/json",
      ...getAuthHeaderFromCookies(),
    },
  });
  if (!res.ok) {
    const resContent = await res.text();
    return { data: null, error: resContent };
  }

  const resJson = await res.json();
  return { data: resJson, error: null };
};

const doPost = async (url, body) => {
  const res = await fetch(url, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      ...getAuthHeaderFromCookies(),
    },
    body: JSON.stringify(body),
  });
  if (!res.ok) {
    const resContent = await res.text();
    return { data: null, error: resContent };
  }

  const resJson = await res.json();
  return { data: resJson, error: null };
};

const doPostAndStreamResponse = async ({ url, reqBody, onDataRecieved }) => {
  const res = await fetch(url, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Connection: "keep-alive",
      ...getAuthHeaderFromCookies(),
    },
    body: JSON.stringify(reqBody),
  });

  if (!res.ok || !res.body) {
    return { error: "Response not ok or body is  empty" };
  }

  const resBody = await res.body;
  const reader = resBody.getReader();
  const decoder = new TextDecoder();

  while (true) {
    const { done, value } = await reader.read();
    if (done) {
      break;
    }
    const chunkString = decoder.decode(value, { stream: true }) || "";
    const jsonStr = chunkString.match("{.*}")?.[0] || "";
    const resultJson = parseJson(jsonStr);

    const predictedText = resultJson?.choices?.[0]?.text;

    onDataRecieved(predictedText);
  }

  return { error: null };
};
