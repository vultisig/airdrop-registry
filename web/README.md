
# Vault File Decode Challenge

Welcome to the Vault File Decode challenge! This challenge will test your skills in React, Next.js, and working with binary data using protobuf.

## Getting Started

1. **Clone the repository**:
   ```
   git clone https://github.com/Ehsan-saradar/Vault-Airdrop
   ```

2. **Install dependencies**:
   ```
   npm install
   ```

3. **Build the project**:
   ```
   npm run build
   ```

4. **Start the development server**:
   ```
   npm run start
   ```

5. **Open the application**:
   Visit [http://localhost:3000](http://localhost:3000) in your browser.

## Challenge Description

You will be implementing a feature to upload and decode a file using protobuf. Follow the steps below to complete the challenge:

### Phase 1

1. **Add UI for File Upload**:
   - Implement a UI that allows the user to upload a file. The UI should support:
     - Uploading a file via a file input.
     - Drag and drop functionality.
     - Removing a file and selecting another one.
   - Only one file can be uploaded at a time.

   - 
### Phase 2
1. **Decode File Using Protobuf**:
   - The uploaded file should be decoded using protobuf. The file contains data of the `VaultContainer` model, which is encoded using protobuf.
   - Decode the file from base64, then use protobuf to decode the byte array into a `VaultContainer` object.

2. **Handle Encrypted Files**:
   - The `VaultContainer` model has an `IsEncrypted` field. If this field is `true`, display a dialog to the user asking for a password. Note: You don't need to decrypt the file at this step, just gather the password.

3. **Extract Vault Name**:
   - Regardless of the `IsEncrypted` value, read the vault data from the `VaultContainer` model. Decode it using base64, and unmarshall it using protobuf.
   - Extract the `Name` field from the vault data.
  
 4. **Display Vault Name**:
   - Redirect the page to `/list-balances`.
   - Set the vault name in a dropdown menu located at the top left of the `/list-balances` page.


## Notes

- Use any protobuf library of your choice for decoding.
- Feel free to add any additional libraries or dependencies as needed.
- Each phase should not take more that 3 hour
- You can use SampleVault.dat as input sample
- You need to use vault.proto and vault_container.proto

## Submission
Please fork this project and create a new branch with your name (e.g., `feature/john-doe`) and commit your changes. Once done, push your branch to the repository and open a pull request.
