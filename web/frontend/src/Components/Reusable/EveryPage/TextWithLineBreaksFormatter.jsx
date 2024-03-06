import React from 'react';

function TextWithLineBreaksFormatter({ text }) {
  // Split the text into paragraphs using '\n\n' as the delimiter
  const paragraphs = text.split('\n\n');

  return (
    <div>
      {paragraphs.map((paragraph, index) => (
        <React.Fragment key={index}>
          {paragraph}
          <br /><br />
        </React.Fragment>
      ))}
    </div>
  );
}

export default TextWithLineBreaksFormatter;
