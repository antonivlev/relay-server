const getAuthHeader = () => {
  const accessToken = document.cookie
    .split(";")
    .find(row => row.startsWith("accessToken"))
    .split("=")?.[1];

  if (!accessToken) {
    return {};
  }

  return {
    Authorization: `Bearer ${accessToken}`,
  };
};

const doPost = async (url, body) => {
  try {
    const res = await fetch(url, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        ...getAuthHeader(),
      },
      body: JSON.stringify(body),
    });
    if (!res.ok) {
      const resContent = await res.text();
      return { data: null, error: resContent };
    }

    const resJson = await res.json();
    return { data: resJson, error: null };
  } catch (error) {
    return { data: null, error };
  }
};
