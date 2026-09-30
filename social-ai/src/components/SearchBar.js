import React, { useState } from "react";
import { Input, Radio } from "antd";
import { SEARCH_KEY } from "../constants";

const { Search } = Input;

function SearchBar({ handleSearch }) {
  const [searchType, setSearchType] = useState(SEARCH_KEY.all);
  const [error, setError] = useState("");

  const changeSearchType = (event) => {
    const selectedType = event.target.value;
    setSearchType(selectedType);
    setError("");

    if (selectedType === SEARCH_KEY.all) {
      handleSearch({ type: SEARCH_KEY.all, keyword: "" });
    }
  };

  const runSearch = (value) => {
    if (searchType !== SEARCH_KEY.all && value.trim() === "") {
      setError("Please input your search keyword.");
      return;
    }

    setError("");
    handleSearch({ type: searchType, keyword: value.trim() });
  };

  return (
    <div className="search-bar">
      <Search
        placeholder="Input search text"
        enterButton="Search"
        size="large"
        onSearch={runSearch}
        disabled={searchType === SEARCH_KEY.all}
      />
      <p className="error-msg">{error}</p>
      <Radio.Group
        onChange={changeSearchType}
        value={searchType}
        className="search-type-group"
      >
        <Radio value={SEARCH_KEY.all}>All</Radio>
        <Radio value={SEARCH_KEY.keyword}>Keyword</Radio>
        <Radio value={SEARCH_KEY.user}>User</Radio>
      </Radio.Group>
    </div>
  );
}

export default SearchBar;
